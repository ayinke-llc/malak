package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type contactHandler struct {
	cfg                config.Config
	contactRepo        malak.ContactRepository
	contactListRepo    malak.ContactListRepository
	contactShareRepo   malak.ContactShareRepository
	referenceGenerator malak.ReferenceGeneratorOperation
}

type createContactRequest struct {
	GenericRequest

	Email     malak.Email `json:"email,omitempty" validate:"'required'"`
	FirstName *string     `json:"first_name,omitempty" validate:"'required'"`

	LastName *string `json:"last_name,omitempty" validate:"'required'"`
}

func (c *createContactRequest) Validate() error {
	if hermes.IsStringEmpty(c.Email.String()) {
		return errors.New("please provide the email address of the contact")
	}

	firstName := strings.TrimSpace(hermes.DeRef(c.FirstName))

	if !hermes.IsStringEmpty(firstName) {
		if len(firstName) > 100 {
			return errors.New("contact's last name must be less than 100")
		}
	}

	c.FirstName = hermes.Ref(firstName)

	lastName := strings.TrimSpace(hermes.DeRef(c.LastName))

	if !hermes.IsStringEmpty(lastName) {
		if len(lastName) > 100 {
			return errors.New("contact's last name must be less than 100")
		}
	}

	c.LastName = hermes.Ref(lastName)

	return nil
}

// @Description Creates a new contact
// @Tags contacts
// @Accept  json
// @Produce  json
// @Param message body createContactRequest true "contact request body"
// @Success 200 {object} fetchContactResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts [post]
func (c *contactHandler) Create(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	user := getUserFromContext(r.Context())

	logger.Debug("creating contact")

	req := new(createContactRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	contact := &malak.Contact{
		Email:       req.Email,
		FirstName:   hermes.DeRef(req.FirstName),
		LastName:    hermes.DeRef(req.LastName),
		Metadata:    make(malak.CustomContactMetadata),
		WorkspaceID: getWorkspaceFromContext(r.Context()).ID,
		Reference:   c.referenceGenerator.Generate(malak.EntityTypeContact),
		// Default to the user who created it
		OwnerID:   user.ID,
		CreatedBy: user.ID,
	}

	err := c.contactRepo.Create(ctx, contact)
	if errors.Is(err, malak.ErrContactExists) {
		return newAPIStatus(http.StatusConflict, err.Error()), StatusFailed
	}

	if err != nil {
		logger.
			Error("an error occurred while storing contact to the database", zap.Error(err))
		return newAPIStatus(
			http.StatusInternalServerError,
			"an error occurred while creating contact"), StatusFailed
	}

	return fetchContactResponse{
		APIStatus: newAPIStatus(http.StatusCreated, "contact was successfully created"),
		Contact:   hermes.DeRef(contact),
	}, StatusSuccess
}

type createContactListRequest struct {
	GenericRequest

	Name string `json:"name,omitempty" validate:"required"`
}

func (c *createContactListRequest) Validate() error {
	p := bluemonday.StrictPolicy()

	c.Name = strings.TrimSpace(p.Sanitize(c.Name))

	if hermes.IsStringEmpty(c.Name) {
		return errors.New("please provide the name of your list")
	}

	if len(c.Name) > 50 {
		return errors.New("your list name cannot be more than 50 characters")
	}

	return nil
}

// @Description Create a new contact list
// @Tags contacts
// @id createContactList
// @Accept  json
// @Produce  json
// @Param message body createContactListRequest true "contact list body"
// @Success 200 {object} fetchContactListResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts/lists [post]
func (c *contactHandler) createContactList(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	user := getUserFromContext(r.Context())

	logger.Debug("creating a new contact list")

	req := new(createContactListRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	list := &malak.ContactList{
		WorkspaceID: getWorkspaceFromContext(r.Context()).ID,
		Reference:   c.referenceGenerator.Generate(malak.EntityTypeList),
		CreatedBy:   user.ID,
		Title:       req.Name,
	}

	if err := c.contactListRepo.Create(ctx, list); err != nil {
		logger.
			Error("an error occurred while storing contact list to the database", zap.Error(err))
		return newAPIStatus(
			http.StatusInternalServerError,
			"an error occurred while creating list"), StatusFailed
	}

	return fetchContactListResponse{
		APIStatus: newAPIStatus(http.StatusCreated, "list was successfully created"),
		List:      hermes.DeRef(list),
	}, StatusSuccess
}

// @Description List all created contact lists
// @Tags contacts
// @id fetchContactLists
// @Accept  json
// @Produce  json
// @Param include_emails query boolean false "show emails inside the list"
// @Success 200 {object} fetchContactListsResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts/lists [get]
func (c *contactHandler) fetchContactLists(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("listing all contact lists in this workspace")

	list, mappings, err := c.contactListRepo.List(ctx, &malak.ContactListOptions{
		WorkspaceID:   getWorkspaceFromContext(ctx).ID,
		IncludeEmails: r.URL.Query().Has("include_emails"),
	})
	if err != nil {
		logger.
			Error("an error occurred while listing contact lists", zap.Error(err))
		return newAPIStatus(
			http.StatusInternalServerError,
			"an error occurred while fetching contact lists"), StatusFailed
	}

	mappingsByListID := make(map[uuid.UUID][]malak.ContactListMappingWithContact)
	for _, mapping := range mappings {
		mappingsByListID[mapping.ListID] = append(mappingsByListID[mapping.ListID], mapping)
	}

	responseLists := []struct {
		List     malak.ContactList                     "json:\"list,omitempty\" validate:\"required\""
		Mappings []malak.ContactListMappingWithContact "json:\"mappings,omitempty\" validate:\"required\""
	}{}

	for _, v := range list {
		responseLists = append(responseLists, struct {
			List     malak.ContactList                     "json:\"list,omitempty\" validate:\"required\""
			Mappings []malak.ContactListMappingWithContact "json:\"mappings,omitempty\" validate:\"required\""
		}{
			List:     v,
			Mappings: mappingsByListID[v.ID],
		})
	}

	return fetchContactListsResponse{
		APIStatus: newAPIStatus(http.StatusOK, "list was successfully retrieved"),
		Lists:     responseLists,
	}, StatusSuccess
}

// @Description Edit a contact list
// @Tags contacts
// @id editContactList
// @Accept  json
// @Produce  json
// @Param message body createContactListRequest true "contact list body"
// @Param reference path string required "list unique reference.. e.g list_"
// @Success 200 {object} fetchContactListResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts/lists/{reference} [put]
func (c *contactHandler) editContactList(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	reference := chi.URLParam(r, "reference")

	logger = logger.With(zap.String("reference", reference))

	logger.Debug("editing contact list")

	req := new(createContactListRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	list, err := c.contactListRepo.Get(ctx, malak.FetchContactListOptions{
		Reference:   malak.Reference(reference),
		WorkspaceID: getWorkspaceFromContext(ctx).ID,
	})
	if errors.Is(err, malak.ErrContactListNotFound) {
		return newAPIStatus(
			http.StatusNotFound, err.Error()), StatusFailed
	}

	if err != nil {
		logger.
			Error("an error occurred while fetching contact list", zap.Error(err))
		return newAPIStatus(
			http.StatusInternalServerError,
			"an error occurred while fetching the contact list"), StatusFailed
	}

	if list.Title != req.Name {
		list.Title = req.Name
		if err := c.contactListRepo.Update(ctx, list); err != nil {
			logger.Error("could not update contact list", zap.Error(err))
			return newAPIStatus(http.StatusInternalServerError, "could not update list"),
				StatusFailed
		}
	}

	return fetchContactListResponse{
		APIStatus: newAPIStatus(http.StatusCreated, "list was successfully created"),
		List:      hermes.DeRef(list),
	}, StatusSuccess
}

// @Description delete a contact list
// @Tags contacts
// @id deleteContactList
// @Accept  json
// @Produce  json
// @Param reference path string required "list unique reference.. e.g list_"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts/lists/{reference} [delete]
func (c *contactHandler) deleteContactList(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	reference := chi.URLParam(r, "reference")

	logger = logger.With(zap.String("reference", reference))

	logger.Debug("deleting contact list")

	list, err := c.contactListRepo.Get(ctx, malak.FetchContactListOptions{
		Reference:   malak.Reference(reference),
		WorkspaceID: getWorkspaceFromContext(ctx).ID,
	})
	if errors.Is(err, malak.ErrContactListNotFound) {
		return newAPIStatus(
			http.StatusNotFound, err.Error()), StatusFailed
	}

	if err != nil {
		logger.
			Error("an error occurred while fetching contact list", zap.Error(err))
		return newAPIStatus(
			http.StatusInternalServerError,
			"an error occurred while fetching the contact list"), StatusFailed
	}

	if err := c.contactListRepo.Delete(ctx, list); err != nil {
		logger.Error("could not delete contact list", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not delete list"),
			StatusFailed
	}

	return newAPIStatus(http.StatusCreated, "list was successfully deleted"), StatusSuccess
}

type addContactToListRequest struct {
	Reference malak.Reference

	GenericRequest
}

func (c *addContactToListRequest) Validate() error {
	if hermes.IsStringEmpty(c.Reference.String()) {
		return errors.New("please provide the reference of the contact")
	}

	return nil
}

// @Description add a new contact to a list
// @Tags contacts
// @id addEmailToContactList
// @Accept  json
// @Produce  json
// @Param message body addContactToListRequest true "contact body"
// @Param reference path string required "list unique reference.. e.g list_"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts/lists/{reference} [post]
func (c *contactHandler) addUserToContactList(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	reference := chi.URLParam(r, "reference")

	logger = logger.With(zap.String("reference", reference))

	logger.Debug("adding a user to a contact list")

	req := new(addContactToListRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	logger = logger.With(zap.String("contact_id", req.Reference.String()),
		zap.String("list_reference", reference))

	contact, err := c.contactRepo.Get(ctx, malak.FetchContactOptions{
		Reference:   req.Reference,
		WorkspaceID: getWorkspaceFromContext(ctx).ID,
	})
	if errors.Is(err, malak.ErrContactNotFound) {
		return newAPIStatus(http.StatusNotFound, err.Error()), StatusFailed
	}

	if err != nil {
		logger.Error("could not fetch contact from the database",
			zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not fetch contact"), StatusFailed
	}

	list, err := c.contactListRepo.Get(ctx, malak.FetchContactListOptions{
		Reference:   malak.Reference(reference),
		WorkspaceID: getWorkspaceFromContext(ctx).ID,
	})
	if errors.Is(err, malak.ErrContactListNotFound) {
		return newAPIStatus(
			http.StatusNotFound, err.Error()), StatusFailed
	}

	if err != nil {
		logger.Error("an error occurred while fetching contact list", zap.Error(err))
		return newAPIStatus(
			http.StatusInternalServerError,
			"an error occurred while fetching the contact list"), StatusFailed
	}

	mapping := &malak.ContactListMapping{
		Reference: c.referenceGenerator.Generate(malak.EntityTypeListEmail),
		ListID:    list.ID,
		ContactID: contact.ID,
		CreatedBy: getUserFromContext(ctx).ID,
	}

	if err := c.contactListRepo.Add(ctx, mapping); err != nil {
		logger.Error("could not add contact list mapping", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not add contact list mapping"),
			StatusFailed
	}

	return newAPIStatus(http.StatusCreated, "list was successfully updated with contact"), StatusSuccess
}

// @Description list your contacts
// @Tags contacts
// @Accept  json
// @Produce  json
// @Param page query int false "Page to query data from. Defaults to 1"
// @Param per_page query int false "Number to items to return. Defaults to 10 items"
// @Success 200 {object} listContactsResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts [get]
func (c *contactHandler) list(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("Listing contacts")

	workspace := getWorkspaceFromContext(r.Context())

	opts := malak.ListContactOptions{
		Paginator:   malak.PaginatorFromRequest(r),
		WorkspaceID: workspace.ID,
	}

	span.SetAttributes(opts.Paginator.OTELAttributes()...)

	contacts, totalCount, err := c.contactRepo.List(ctx, opts)
	if err != nil {
		logger.Error("could not list contacts",
			zap.Error(err))

		return newAPIStatus(
			http.StatusInternalServerError,
			"could not list contacts"), StatusFailed
	}

	return listContactsResponse{
		APIStatus: newAPIStatus(http.StatusCreated, "contact listed successfully"),
		Contacts:  contacts,
		Meta: meta{
			Paging: pagingInfo{
				PerPage: opts.Paginator.PerPage,
				Page:    opts.Paginator.Page,
				Total:   totalCount,
			},
		},
	}, StatusSuccess
}

// @Description fetch a contact by reference
// @Tags contacts
// @Accept  json
// @Produce  json
// @Param reference path string required "contact unique reference.. e.g contact_"
// @Success 200 {object} fetchDetailedContactResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts/{reference} [get]
func (c *contactHandler) fetchContact(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	workspace := getWorkspaceFromContext(r.Context())

	reference := chi.URLParam(r, "reference")

	logger = logger.With(zap.String("reference", reference))

	logger.Debug("fetching a single contact")

	contact, err := c.contactRepo.Get(ctx, malak.FetchContactOptions{
		WorkspaceID: workspace.ID,
		Reference:   malak.Reference(reference),
	})
	if err != nil {
		logger.Error("could not fetch contact",
			zap.Error(err))

		var status = http.StatusInternalServerError
		var msg = "could not fetch contact"

		if errors.Is(err, malak.ErrContactNotFound) {
			status = http.StatusNotFound
			msg = "contact does not exists"
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	var g errgroup.Group
	var sharedItems []malak.ContactShareItem

	// errgroup because of analytics in the future
	g.Go(func() error {
		var err error

		sharedItems, err = c.contactShareRepo.All(ctx, contact)
		return err
	})

	if err := g.Wait(); err != nil {
		logger.Error("could not fetch contact details", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not fetch contact details"),
			StatusFailed
	}

	return fetchDetailedContactResponse{
		APIStatus:   newAPIStatus(http.StatusOK, "contact was retrieved"),
		Contact:     hermes.DeRef(contact),
		SharedItems: sharedItems,
	}, StatusSuccess
}

// @Description delete a contact
// @Tags contacts
// @id deleteContact
// @Accept  json
// @Produce  json
// @Param reference path string required "contact unique reference.. e.g contact_"
// @Success 200 {object} APIStatus
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts/{reference} [delete]
func (c *contactHandler) deleteContact(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	reference := chi.URLParam(r, "reference")

	logger = logger.With(zap.String("reference", reference))

	logger.Debug("deleting a single contact")

	workspace := getWorkspaceFromContext(r.Context())

	contact, err := c.contactRepo.Get(ctx, malak.FetchContactOptions{
		WorkspaceID: workspace.ID,
		Reference:   malak.Reference(reference),
	})
	if err != nil {
		logger.Error("could not fetch contact",
			zap.Error(err))

		var status = http.StatusInternalServerError
		var msg = "could not fetch contact"

		if errors.Is(err, malak.ErrContactNotFound) {
			status = http.StatusNotFound
			msg = "contact does not exists"
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	if err := c.contactRepo.Delete(ctx, contact); err != nil {
		logger.Error("could not delete contact list", zap.Error(err))
		return newAPIStatus(http.StatusInternalServerError, "could not delete list"),
			StatusFailed
	}

	return newAPIStatus(http.StatusOK, "contact was successfully deleted"), StatusSuccess
}

type editContactRequest struct {
	GenericRequest

	FirstName string `json:"first_name,omitempty" validate:"'required'"`
	LastName  string `json:"last_name,omitempty" validate:"'required'"`
	Company   string `json:"company,omitempty" validate:"required"`
	Address   string `json:"address,omitempty" validate:"required"`
	Notes     string `json:"notes,omitempty" validate:"required"`
}

func (c *editContactRequest) Validate() error {

	c.FirstName = strings.TrimSpace(c.FirstName)

	if !hermes.IsStringEmpty(c.FirstName) {
		if len(c.FirstName) > 50 {
			return errors.New("contact's first name must be less than 50")
		}

		if len(c.FirstName) < 5 {
			return errors.New("contact's first name must be greater than 5")
		}
	}

	c.LastName = strings.TrimSpace(c.LastName)

	if !hermes.IsStringEmpty(c.LastName) {
		if len(c.LastName) > 50 {
			return errors.New("contact's last name must be less than 50")
		}

		if len(c.LastName) < 5 {
			return errors.New("contact's last name must be greater than 5")
		}
	}

	c.Company = strings.TrimSpace(c.Company)

	if !hermes.IsStringEmpty(c.Company) {
		if len(c.Company) > 100 {
			return errors.New("company's name must be less than 100")
		}

		if len(c.Company) < 5 {
			return errors.New("company's name must be greater than 5 characters")
		}
	}

	c.Address = strings.TrimSpace(c.Address)

	if !hermes.IsStringEmpty(c.Address) {
		if len(c.Address) > 225 {
			return errors.New("company's address must be less than 225")
		}

		if len(c.Address) < 5 {
			return errors.New("company's address must be greater than 5 characters")
		}
	}

	c.Notes = strings.TrimSpace(c.Notes)

	if !hermes.IsStringEmpty(c.Notes) {
		if len(c.Notes) > 2000 {
			return errors.New("company's notes must be less than 2000")
		}

		if len(c.Notes) < 5 {
			return errors.New("company's notes must be greater than 5 characters")
		}
	}

	return nil
}

// @Description edit a contact
// @Tags contacts
// @Accept  json
// @Produce  json
// @Param message body editContactRequest true "contact request body"
// @Param reference path string required "contact unique reference.. e.g contact_"
// @Success 200 {object} fetchContactResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts/{reference} [put]
func (c *contactHandler) editContact(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	workspace := getWorkspaceFromContext(ctx)

	reference := chi.URLParam(r, "reference")

	logger.Debug("editing contact")

	req := new(editContactRequest)

	if err := render.Bind(r, req); err != nil {
		return newAPIStatus(http.StatusBadRequest, "invalid request body"), StatusFailed
	}

	if err := req.Validate(); err != nil {
		return newAPIStatus(http.StatusBadRequest, err.Error()), StatusFailed
	}

	contact, err := c.contactRepo.Get(ctx, malak.FetchContactOptions{
		WorkspaceID: workspace.ID,
		Reference:   malak.Reference(reference),
	})
	if err != nil {
		logger.Error("could not fetch contact",
			zap.Error(err))

		var status = http.StatusInternalServerError
		var msg = "could not fetch contact"

		if errors.Is(err, malak.ErrContactNotFound) {
			status = http.StatusNotFound
			msg = "contact does not exists"
		}

		return newAPIStatus(status, msg), StatusFailed
	}

	if req.Address != contact.City {
		contact.City = req.Address
	}

	if req.Company != contact.Company {
		contact.Company = req.Company
	}

	if req.FirstName != contact.FirstName {
		contact.FirstName = req.FirstName
	}

	if req.LastName != contact.LastName {
		contact.LastName = req.LastName
	}

	if req.Notes != contact.Notes {
		contact.Notes = req.Notes
	}

	if err := c.contactRepo.Update(ctx, contact); err != nil {
		logger.
			Error("an error occurred while updating contact", zap.Error(err))
		return newAPIStatus(
			http.StatusInternalServerError,
			"an error occurred while updating contact"), StatusFailed
	}

	return fetchContactResponse{
		APIStatus: newAPIStatus(http.StatusOK, "contact was successfully updated"),
		Contact:   hermes.DeRef(contact),
	}, StatusSuccess
}

// @Description Search contacts
// @Tags contacts
// @Accept  json
// @Produce  json
// @Param search query string true "search term"
// @Success 200 {object} listContactsResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts/search [get]
func (c *contactHandler) search(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("searching contacts")

	workspace := getWorkspaceFromContext(r.Context())
	searchValue := r.URL.Query().Get("search")

	if searchValue == "" {
		return newAPIStatus(http.StatusBadRequest, "search parameter is required"), StatusFailed
	}

	contacts, err := c.contactRepo.Search(ctx, malak.SearchContactOptions{
		WorkspaceID: workspace.ID,
		SearchValue: searchValue,
	})
	if err != nil {
		logger.Error("could not search contacts",
			zap.Error(err))

		return newAPIStatus(
			http.StatusInternalServerError,
			"could not search contacts"), StatusFailed
	}

	return listContactsResponse{
		APIStatus: newAPIStatus(http.StatusOK, "contacts searched successfully"),
		Contacts:  contacts,
		Meta: meta{
			Paging: pagingInfo{
				PerPage: int64(len(contacts)),
				Page:    1,
				Total:   int64(len(contacts)),
			},
		},
	}, StatusSuccess
}

// @Description List all contacts
// @Tags contacts
// @id listAllContacts
// @Accept  json
// @Produce  json
// @Success 200 {object} listContactsResponse
// @Failure 400 {object} APIStatus
// @Failure 401 {object} APIStatus
// @Failure 404 {object} APIStatus
// @Failure 500 {object} APIStatus
// @Router /contacts/all [get]
func (c *contactHandler) listContacts(
	ctx context.Context,
	span trace.Span,
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) (render.Renderer, Status) {

	logger.Debug("listing all contact in this workspace")

	// TODO(adelowo): Keep track of otel for this to know if we should cut this out and use
	// search instead
	contacts, err := c.contactRepo.All(ctx, getWorkspaceFromContext(ctx).ID)
	if err != nil {
		logger.
			Error("an error occurred while listing contacts", zap.Error(err))
		return newAPIStatus(
			http.StatusInternalServerError,
			"an error occurred while fetching contacts"), StatusFailed
	}

	return listContactsResponse{
		APIStatus: newAPIStatus(http.StatusOK, "fetched all contacts"),
		Contacts:  contacts,
	}, StatusSuccess
}
