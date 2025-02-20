/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface MalakBillingPreferences {
  finance_email?: string;
}

export interface MalakBlock {
  children?: MalakBlock[];
  content?: any;
  id?: string;
  props?: Record<string, any>;
  type?: string;
}

export interface MalakCommunicationPreferences {
  enable_marketing?: boolean;
  enable_product_updates?: boolean;
}

export interface MalakContact {
  /** Legacy lmao. should be address but migrations bit ugh :)) */
  city?: string;
  company?: string;
  created_at?: string;
  /** User who added/created this contact */
  created_by?: string;
  email?: string;
  first_name?: string;
  id?: string;
  last_name?: string;
  lists?: MalakContactListMapping[];
  metadata?: MalakCustomContactMetadata;
  notes?: string;
  /**
   * User who owns the contact.
   * Does not mean who added the contact but who chases
   * or follows up officially with the contact
   */
  owner_id?: string;
  phone?: string;
  reference?: string;
  updated_at?: string;
  workspace_id?: string;
}

export interface MalakContactList {
  created_at?: string;
  created_by?: string;
  id?: string;
  reference?: string;
  title?: string;
  updated_at?: string;
  workspace_id?: string;
}

export interface MalakContactListMapping {
  contact_id?: string;
  created_at?: string;
  created_by?: string;
  id?: string;
  list?: MalakContactList;
  list_id?: string;
  reference?: string;
  updated_at?: string;
}

export interface MalakContactListMappingWithContact {
  contact_id?: string;
  /** Contact fields */
  email?: string;
  id?: string;
  list_id?: string;
  reference?: string;
}

export interface MalakContactShareItem {
  contact_id?: string;
  created_at?: string;
  id?: string;
  item_id?: string;
  item_reference?: string;
  item_type?: MalakContactShareItemType;
  reference?: string;
  shared_at?: string;
  shared_by?: string;
  title?: string;
  updated_at?: string;
}

export enum MalakContactShareItemType {
  ContactShareItemTypeUpdate = "update",
  ContactShareItemTypeDashboard = "dashboard",
  ContactShareItemTypeDeck = "deck",
}

export type MalakCustomContactMetadata = Record<string, string>;

export interface MalakDashboard {
  chart_count?: number;
  created_at?: string;
  description?: string;
  id?: string;
  reference?: string;
  title?: string;
  updated_at?: string;
  workspace_id?: string;
}

export interface MalakDeck {
  created_at?: string;
  created_by?: string;
  deck_size?: number;
  id?: string;
  is_archived?: boolean;
  is_pinned?: boolean;
  object_key?: string;
  preferences?: MalakDeckPreference;
  reference?: string;
  short_link?: string;
  title?: string;
  updated_at?: string;
  workspace_id?: string;
}

export interface MalakDeckPreference {
  created_at?: string;
  created_by?: string;
  deck_id?: string;
  enable_downloading?: boolean;
  expires_at?: string;
  id?: string;
  password?: MalakPasswordDeckPreferences;
  reference?: string;
  require_email?: boolean;
  updated_at?: string;
  workspace_id?: string;
}

export interface MalakIntegration {
  created_at?: string;
  description?: string;
  id?: string;
  integration_name?: string;
  integration_type?: MalakIntegrationType;
  is_enabled?: boolean;
  logo_url?: string;
  metadata?: MalakIntegrationMetadata;
  reference?: string;
  updated_at?: string;
}

export interface MalakIntegrationChart {
  created_at?: string;
  id?: string;
  internal_name?: MalakIntegrationChartInternalNameType;
  metadata?: MalakIntegrationChartMetadata;
  reference?: string;
  updated_at?: string;
  user_facing_name?: string;
  workspace_id?: string;
  workspace_integration_id?: string;
}

export enum MalakIntegrationChartInternalNameType {
  IntegrationChartInternalNameTypeMercuryAccount = "mercury_account",
  IntegrationChartInternalNameTypeMercuryAccountTransaction = "mercury_account_transaction",
  IntegrationChartInternalNameTypeBrexAccount = "brex_account",
  IntegrationChartInternalNameTypeBrexAccountTransaction = "brex_account_transaction",
}

export interface MalakIntegrationChartMetadata {
  provider_id?: string;
}

export interface MalakIntegrationMetadata {
  endpoint?: string;
}

export enum MalakIntegrationType {
  IntegrationTypeOauth2 = "oauth2",
  IntegrationTypeApiKey = "api_key",
}

export interface MalakPasswordDeckPreferences {
  enabled?: boolean;
  password?: string;
}

export interface MalakPlan {
  /** Defaults to zero */
  amount?: number;
  created_at?: string;
  /** Stripe default price id. Again not needed if not using Stripe */
  default_price_id?: string;
  id?: string;
  /**
   * IsDefault if this is the default plan for the user to get signed up to
   * on sign up
   *
   * Better to keep this here than to use config
   */
  is_default?: boolean;
  metadata?: MalakPlanMetadata;
  plan_name?: string;
  /**
   * Can use a fake id really
   * As this only matters if you turn on Stripe
   */
  reference?: string;
  updated_at?: string;
}

export interface MalakPlanMetadata {
  dashboard?: {
    embed_dashboard?: boolean;
    share_dashboard_via_link?: boolean;
  };
  data_room?: {
    share_via_link?: boolean;
    size?: number;
  };
  deck?: {
    auto_terminate_link?: boolean;
    custom_domain?: boolean;
  };
  integrations?: {
    available_for_use?: number;
  };
  team?: {
    size?: number;
  };
  updates?: {
    custom_domain?: boolean;
    max_recipients?: number;
  };
}

export interface MalakPreference {
  billing?: MalakBillingPreferences;
  communication?: MalakCommunicationPreferences;
  created_at?: string;
  id?: string;
  updated_at?: string;
  workspace_id?: string;
}

export interface MalakPublicDeck {
  created_at?: string;
  deck_size?: number;
  is_archived?: boolean;
  object_link?: string;
  preferences?: MalakDeckPreference;
  reference?: string;
  short_link?: string;
  title?: string;
  updated_at?: string;
  workspace_id?: string;
}

export enum MalakRecipientStatus {
  RecipientStatusPending = "pending",
  RecipientStatusSent = "sent",
  RecipientStatusFailed = "failed",
}

export enum MalakRole {
  RoleAdmin = "admin",
  RoleMember = "member",
  RoleBilling = "billing",
  RoleInvestor = "investor",
  RoleGuest = "guest",
}

export interface MalakUpdate {
  content?: MalakBlock[];
  created_at?: string;
  created_by?: string;
  id?: string;
  /** If this update is pinned */
  is_pinned?: boolean;
  metadata?: MalakUpdateMetadata;
  reference?: string;
  sent_at?: string;
  sent_by?: string;
  status?: MalakUpdateStatus;
  title?: string;
  updated_at?: string;
  workspace_id?: string;
}

export type MalakUpdateMetadata = object;

export interface MalakUpdateRecipient {
  contact?: MalakContact;
  contact_id?: string;
  created_at?: string;
  id?: string;
  reference?: string;
  schedule_id?: string;
  status?: MalakRecipientStatus;
  update_id?: string;
  update_recipient_stat?: MalakUpdateRecipientStat;
  updated_at?: string;
}

export interface MalakUpdateRecipientStat {
  created_at?: string;
  has_reaction?: boolean;
  id?: string;
  is_bounced?: boolean;
  is_delivered?: boolean;
  last_opened_at?: string;
  recipient?: MalakUpdateRecipient;
  recipient_id?: string;
  reference?: string;
  updated_at?: string;
}

export interface MalakUpdateStat {
  created_at?: string;
  id?: string;
  reference?: string;
  total_clicks?: number;
  total_opens?: number;
  total_reactions?: number;
  total_sent?: number;
  unique_opens?: number;
  update_id?: string;
  updated_at?: string;
}

export enum MalakUpdateStatus {
  UpdateStatusDraft = "draft",
  UpdateStatusSent = "sent",
}

export interface MalakUser {
  created_at?: string;
  email?: string;
  full_name?: string;
  id?: string;
  metadata?: MalakUserMetadata;
  roles?: MalakUserRole[];
  updated_at?: string;
}

export interface MalakUserMetadata {
  /**
   * Used to keep track of the last used workspace
   * In the instance of multiple workspaces
   * So when next the user logs in, we remember and take them to the
   * right place rather than always a list of all their workspaces and they
   * have to select one
   */
  current_workspace?: string;
}

export interface MalakUserRole {
  created_at?: string;
  id?: string;
  role?: MalakRole;
  updated_at?: string;
  user_id?: string;
  workspace_id?: string;
}

export interface MalakWorkspace {
  created_at?: string;
  id?: string;
  is_subscription_active?: boolean;
  logo_url?: string;
  metadata?: MalakWorkspaceMetadata;
  plan?: MalakPlan;
  plan_id?: string;
  reference?: string;
  /**
   * Not required
   * Dummy values work really if not using stripe
   */
  stripe_customer_id?: string;
  subscription_id?: string;
  timezone?: string;
  updated_at?: string;
  website?: string;
  workspace_name?: string;
}

export interface MalakWorkspaceIntegration {
  created_at?: string;
  id?: string;
  integration?: MalakIntegration;
  integration_id?: string;
  /** IsActive determines if the connection to the integration has been tested and works */
  is_active?: boolean;
  /** IsEnabled - this integration is enabled and data can be fetched */
  is_enabled?: boolean;
  metadata?: MalakWorkspaceIntegrationMetadata;
  reference?: string;
  updated_at?: string;
  workspace_id?: string;
}

export interface MalakWorkspaceIntegrationMetadata {
  access_token?: string;
  last_fetched_at?: string;
}

export type MalakWorkspaceMetadata = object;

export interface ServerAPIStatus {
  message: string;
}

export interface ServerAddContactToListRequest {
  reference?: string;
}

export interface ServerAuthenticateUserRequest {
  code: string;
}

export interface ServerContentUpdateRequest {
  title: string;
  update: MalakBlock[];
}

export interface ServerCreateContactListRequest {
  name: string;
}

export interface ServerCreateContactRequest {
  email?: string;
  first_name?: string;
  last_name?: string;
}

export interface ServerCreateDashboardRequest {
  description?: string;
  title?: string;
}

export interface ServerCreateDeckRequest {
  deck_url?: string;
  title?: string;
}

export interface ServerCreateUpdateContent {
  title: string;
}

export interface ServerCreateWorkspaceRequest {
  name: string;
}

export interface ServerCreatedUpdateResponse {
  message: string;
  update: MalakUpdate;
}

export interface ServerCreatedUserResponse {
  current_workspace?: MalakWorkspace;
  message: string;
  token: string;
  user: MalakUser;
  workspaces: MalakWorkspace[];
}

export interface ServerEditContactRequest {
  address: string;
  company: string;
  first_name?: string;
  last_name?: string;
  notes: string;
}

export interface ServerFetchBillingPortalResponse {
  link: string;
  message: string;
}

export interface ServerFetchContactListResponse {
  list: MalakContactList;
  message: string;
}

export interface ServerFetchContactListsResponse {
  lists: {
    list: MalakContactList;
    mappings: MalakContactListMappingWithContact[];
  }[];
  message: string;
}

export interface ServerFetchContactResponse {
  contact: MalakContact;
  message: string;
}

export interface ServerFetchDashboardResponse {
  dashboard: MalakDashboard;
  message: string;
}

export interface ServerFetchDeckResponse {
  deck: MalakDeck;
  message: string;
}

export interface ServerFetchDecksResponse {
  decks: MalakDeck[];
  message: string;
}

export interface ServerFetchDetailedContactResponse {
  contact: MalakContact;
  message: string;
  shared_items: MalakContactShareItem[];
}

export interface ServerFetchPublicDeckResponse {
  deck: MalakPublicDeck;
  message: string;
}

export interface ServerFetchUpdateAnalyticsResponse {
  message: string;
  recipients: MalakUpdateRecipient[];
  update: MalakUpdateStat;
}

export interface ServerFetchUpdateReponse {
  message: string;
  update: MalakUpdate;
}

export interface ServerFetchWorkspaceResponse {
  message: string;
  workspace: MalakWorkspace;
}

export interface ServerListContactsResponse {
  contacts: MalakContact[];
  message: string;
  meta: ServerMeta;
}

export interface ServerListDashboardResponse {
  dashboards: MalakDashboard[];
  message: string;
  meta: ServerMeta;
}

export interface ServerListIntegrationChartsResponse {
  charts: MalakIntegrationChart[];
  message: string;
}

export interface ServerListIntegrationResponse {
  integrations: MalakWorkspaceIntegration[];
  message: string;
}

export interface ServerListUpdateResponse {
  message: string;
  meta: ServerMeta;
  updates: MalakUpdate[];
}

export interface ServerMeta {
  paging: ServerPagingInfo;
}

export interface ServerPagingInfo {
  page: number;
  per_page: number;
  total: number;
}

export interface ServerPreferenceResponse {
  message: string;
  preferences: MalakPreference;
}

export interface ServerPreviewUpdateRequest {
  email: string;
}

export interface ServerSendUpdateRequest {
  emails?: string[];
  send_at?: number;
}

export interface ServerTestAPIIntegrationRequest {
  api_key: string;
}

export interface ServerUpdateDeckPreferencesRequest {
  enable_downloading?: boolean;
  password_protection?: {
    enabled?: boolean;
    value?: string;
  };
  require_email?: boolean;
}

export interface ServerUpdatePreferencesRequest {
  preferences: {
    billing: MalakBillingPreferences;
    newsletter: MalakCommunicationPreferences;
  };
}

export interface ServerUpdateWorkspaceRequest {
  logo?: string;
  timezone?: string;
  website?: string;
  workspace_name?: string;
}

export interface ServerUploadImageResponse {
  message: string;
  url: string;
}

import type { AxiosInstance, AxiosRequestConfig, AxiosResponse, HeadersDefaults, ResponseType } from "axios";
import axios from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({ securityWorker, secure, format, ...axiosConfig }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "http://localhost:5300/v1" });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected mergeRequestParams(params1: AxiosRequestConfig, params2?: AxiosRequestConfig): AxiosRequestConfig {
    const method = params1.method || (params2 && params2.method);

    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...((method && this.instance.defaults.headers[method.toLowerCase() as keyof HeadersDefaults]) || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected stringifyFormItem(formItem: unknown) {
    if (typeof formItem === "object" && formItem !== null) {
      return JSON.stringify(formItem);
    } else {
      return `${formItem}`;
    }
  }

  protected createFormData(input: Record<string, unknown>): FormData {
    if (input instanceof FormData) {
      return input;
    }
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      const propertyContent: any[] = property instanceof Array ? property : [property];

      for (const formItem of propertyContent) {
        const isFileType = formItem instanceof Blob || formItem instanceof File;
        formData.append(key, isFileType ? formItem : this.stringifyFormItem(formItem));
      }

      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = format || this.format || undefined;

    if (type === ContentType.FormData && body && body !== null && typeof body === "object") {
      body = this.createFormData(body as Record<string, unknown>);
    }

    if (type === ContentType.Text && body && body !== null && typeof body !== "string") {
      body = JSON.stringify(body);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(requestParams.headers || {}),
        ...(type ? { "Content-Type": type } : {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title Malak's API documentation
 * @version 0.1.0
 * @baseUrl http://localhost:5300/v1
 * @contact Ayinke Ventures <lanre@ayinke.ventures>
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  auth = {
    /**
     * No description
     *
     * @tags auth
     * @name ConnectCreate
     * @summary Sign in with a social login provider
     * @request POST:/auth/connect/{provider}
     */
    connectCreate: (provider: string, data: ServerAuthenticateUserRequest, params: RequestParams = {}) =>
      this.request<ServerCreatedUserResponse, ServerAPIStatus>({
        path: `/auth/connect/${provider}`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  contacts = {
    /**
     * No description
     *
     * @tags contacts
     * @name ContactsList
     * @summary list your contacts
     * @request GET:/contacts
     */
    contactsList: (
      query?: {
        /** Page to query data from. Defaults to 1 */
        page?: number;
        /** Number to items to return. Defaults to 10 items */
        per_page?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<ServerListContactsResponse, ServerAPIStatus>({
        path: `/contacts`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags contacts
     * @name ContactsCreate
     * @summary Creates a new contact
     * @request POST:/contacts
     */
    contactsCreate: (data: ServerCreateContactRequest, params: RequestParams = {}) =>
      this.request<ServerFetchContactResponse, ServerAPIStatus>({
        path: `/contacts`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags contacts
     * @name DeleteContact
     * @summary delete a contact
     * @request DELETE:/contacts/{reference}
     */
    deleteContact: (reference: string, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/contacts/${reference}`,
        method: "DELETE",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags contacts
     * @name ContactsDetail
     * @summary fetch a contact by reference
     * @request GET:/contacts/{reference}
     */
    contactsDetail: (reference: string, params: RequestParams = {}) =>
      this.request<ServerFetchDetailedContactResponse, ServerAPIStatus>({
        path: `/contacts/${reference}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags contacts
     * @name ContactsUpdate
     * @summary edit a contact
     * @request PUT:/contacts/{reference}
     */
    contactsUpdate: (reference: string, data: ServerEditContactRequest, params: RequestParams = {}) =>
      this.request<ServerFetchContactResponse, ServerAPIStatus>({
        path: `/contacts/${reference}`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags contacts
     * @name FetchContactLists
     * @summary List all created contact lists
     * @request GET:/contacts/lists
     */
    fetchContactLists: (
      query?: {
        /** show emails inside the list */
        include_emails?: boolean;
      },
      params: RequestParams = {},
    ) =>
      this.request<ServerFetchContactListsResponse, ServerAPIStatus>({
        path: `/contacts/lists`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags contacts
     * @name CreateContactList
     * @summary Create a new contact list
     * @request POST:/contacts/lists
     */
    createContactList: (data: ServerCreateContactListRequest, params: RequestParams = {}) =>
      this.request<ServerFetchContactListResponse, ServerAPIStatus>({
        path: `/contacts/lists`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags contacts
     * @name DeleteContactList
     * @summary delete a contact list
     * @request DELETE:/contacts/lists/{reference}
     */
    deleteContactList: (reference: string, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/contacts/lists/${reference}`,
        method: "DELETE",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags contacts
     * @name AddEmailToContactList
     * @summary add a new contact to a list
     * @request POST:/contacts/lists/{reference}
     */
    addEmailToContactList: (reference: string, data: ServerAddContactToListRequest, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/contacts/lists/${reference}`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags contacts
     * @name EditContactList
     * @summary Edit a contact list
     * @request PUT:/contacts/lists/{reference}
     */
    editContactList: (reference: string, data: ServerCreateContactListRequest, params: RequestParams = {}) =>
      this.request<ServerFetchContactListResponse, ServerAPIStatus>({
        path: `/contacts/lists/${reference}`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  dashboards = {
    /**
     * No description
     *
     * @tags dashboards
     * @name DashboardsList
     * @summary List dashboards
     * @request GET:/dashboards
     */
    dashboardsList: (
      query?: {
        /** Page to query data from. Defaults to 1 */
        page?: number;
        /** Number to items to return. Defaults to 10 items */
        per_page?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<ServerListDashboardResponse, ServerAPIStatus>({
        path: `/dashboards`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags dashboards
     * @name DashboardsCreate
     * @summary create a new dashboard
     * @request POST:/dashboards
     */
    dashboardsCreate: (data: ServerCreateDashboardRequest, params: RequestParams = {}) =>
      this.request<ServerFetchDashboardResponse, ServerAPIStatus>({
        path: `/dashboards`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags dashboards
     * @name ChartsList
     * @summary List charts
     * @request GET:/dashboards/charts
     */
    chartsList: (params: RequestParams = {}) =>
      this.request<ServerListIntegrationChartsResponse, ServerAPIStatus>({
        path: `/dashboards/charts`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
  decks = {
    /**
     * No description
     *
     * @tags decks
     * @name DecksList
     * @summary list all decks. No pagination
     * @request GET:/decks
     */
    decksList: (params: RequestParams = {}) =>
      this.request<ServerFetchDecksResponse, ServerAPIStatus>({
        path: `/decks`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags decks
     * @name DecksCreate
     * @summary Creates a new deck
     * @request POST:/decks
     */
    decksCreate: (data: ServerCreateDeckRequest, params: RequestParams = {}) =>
      this.request<ServerFetchDeckResponse, ServerAPIStatus>({
        path: `/decks`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags decks
     * @name DecksDelete
     * @summary delete a deck
     * @request DELETE:/decks/{reference}
     */
    decksDelete: (reference: string, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/decks/${reference}`,
        method: "DELETE",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags decks
     * @name DecksDetail
     * @summary fetch a deck
     * @request GET:/decks/{reference}
     */
    decksDetail: (reference: string, params: RequestParams = {}) =>
      this.request<ServerFetchDeckResponse, ServerAPIStatus>({
        path: `/decks/${reference}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags decks
     * @name ToggleArchive
     * @summary toggle archive status of a deck
     * @request POST:/decks/{reference}/archive
     */
    toggleArchive: (reference: string, params: RequestParams = {}) =>
      this.request<ServerFetchDeckResponse, ServerAPIStatus>({
        path: `/decks/${reference}/archive`,
        method: "POST",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags decks
     * @name TogglePin
     * @summary toggle pinned status of a deck
     * @request POST:/decks/{reference}/pin
     */
    togglePin: (reference: string, params: RequestParams = {}) =>
      this.request<ServerFetchDeckResponse, ServerAPIStatus>({
        path: `/decks/${reference}/pin`,
        method: "POST",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags decks
     * @name PreferencesUpdate
     * @summary update a deck preferences
     * @request PUT:/decks/{reference}/preferences
     */
    preferencesUpdate: (reference: string, data: ServerUpdateDeckPreferencesRequest, params: RequestParams = {}) =>
      this.request<ServerFetchDeckResponse, ServerAPIStatus>({
        path: `/decks/${reference}/preferences`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  public = {
    /**
     * No description
     *
     * @tags decks-viewer
     * @name DecksDetail
     * @summary public api to fetch a deck
     * @request GET:/public/decks/{reference}
     */
    decksDetail: (reference: string, params: RequestParams = {}) =>
      this.request<ServerFetchPublicDeckResponse, ServerAPIStatus>({
        path: `/public/decks/${reference}`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
  updates = {
    /**
     * No description
     *
     * @tags updates
     * @name ReactPost
     * @summary Fetch a specific update
     * @request GET:/updates/react
     */
    reactPost: (
      query: {
        /** provider type */
        provider: string;
        /** email id */
        email_id: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/updates/react`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),
  };
  uploads = {
    /**
     * No description
     *
     * @tags decks
     * @name UploadDeck
     * @summary Upload a deck
     * @request POST:/uploads/decks
     */
    uploadDeck: (
      data: {
        /**
         * image body
         * @format binary
         */
        image_body: File;
      },
      params: RequestParams = {},
    ) =>
      this.request<ServerUploadImageResponse, ServerAPIStatus>({
        path: `/uploads/decks`,
        method: "POST",
        body: data,
        type: ContentType.FormData,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags images
     * @name UploadImage
     * @summary Upload an image
     * @request POST:/uploads/images
     */
    uploadImage: (
      data: {
        /**
         * image body
         * @format binary
         */
        image_body: File;
      },
      params: RequestParams = {},
    ) =>
      this.request<ServerUploadImageResponse, ServerAPIStatus>({
        path: `/uploads/images`,
        method: "POST",
        body: data,
        type: ContentType.FormData,
        format: "json",
        ...params,
      }),
  };
  user = {
    /**
     * No description
     *
     * @tags user
     * @name UserList
     * @summary Fetch current user. This api should also double as a token validation api
     * @request GET:/user
     */
    userList: (params: RequestParams = {}) =>
      this.request<ServerCreatedUserResponse, ServerAPIStatus>({
        path: `/user`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
  workspaces = {
    /**
     * No description
     *
     * @tags workspace
     * @name WorkspacesPartialUpdate
     * @summary update workspace details
     * @request PATCH:/workspaces
     */
    workspacesPartialUpdate: (data: ServerUpdateWorkspaceRequest, params: RequestParams = {}) =>
      this.request<ServerFetchWorkspaceResponse, ServerAPIStatus>({
        path: `/workspaces`,
        method: "PATCH",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags workspace
     * @name WorkspacesCreate
     * @summary Create a new workspace
     * @request POST:/workspaces
     */
    workspacesCreate: (data: ServerCreateWorkspaceRequest, params: RequestParams = {}) =>
      this.request<ServerFetchWorkspaceResponse, ServerAPIStatus>({
        path: `/workspaces`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags billing
     * @name BillingCreate
     * @summary get billing portal
     * @request POST:/workspaces/billing
     */
    billingCreate: (params: RequestParams = {}) =>
      this.request<ServerFetchBillingPortalResponse, ServerAPIStatus>({
        path: `/workspaces/billing`,
        method: "POST",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags integrations
     * @name IntegrationsList
     * @summary fetch workspace preferences
     * @request GET:/workspaces/integrations
     */
    integrationsList: (params: RequestParams = {}) =>
      this.request<ServerListIntegrationResponse, ServerAPIStatus>({
        path: `/workspaces/integrations`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags integrations
     * @name IntegrationsDelete
     * @summary disable integration
     * @request DELETE:/workspaces/integrations/{reference}
     */
    integrationsDelete: (reference: string, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/workspaces/integrations/${reference}`,
        method: "DELETE",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags integrations
     * @name IntegrationsCreate
     * @summary enable integration
     * @request POST:/workspaces/integrations/{reference}
     */
    integrationsCreate: (reference: string, data: ServerTestAPIIntegrationRequest, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/workspaces/integrations/${reference}`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags integrations
     * @name IntegrationsUpdate
     * @summary update integration api key
     * @request PUT:/workspaces/integrations/{reference}
     */
    integrationsUpdate: (reference: string, data: ServerTestAPIIntegrationRequest, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/workspaces/integrations/${reference}`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags integrations
     * @name IntegrationsPingCreate
     * @summary test an api key is valid and can reach the integration
     * @request POST:/workspaces/integrations/{reference}/ping
     */
    integrationsPingCreate: (reference: string, data: ServerTestAPIIntegrationRequest, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/workspaces/integrations/${reference}/ping`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags workspace
     * @name PreferencesList
     * @summary fetch workspace preferences
     * @request GET:/workspaces/preferences
     */
    preferencesList: (params: RequestParams = {}) =>
      this.request<ServerPreferenceResponse, ServerAPIStatus>({
        path: `/workspaces/preferences`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags workspace
     * @name PreferencesUpdate
     * @summary update workspace preferences
     * @request PUT:/workspaces/preferences
     */
    preferencesUpdate: (data: ServerUpdatePreferencesRequest, params: RequestParams = {}) =>
      this.request<ServerPreferenceResponse, ServerAPIStatus>({
        path: `/workspaces/preferences`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags workspace
     * @name Switchworkspace
     * @summary Switch current workspace
     * @request POST:/workspaces/switch/{reference}
     */
    switchworkspace: (reference: string, params: RequestParams = {}) =>
      this.request<ServerFetchWorkspaceResponse, ServerAPIStatus>({
        path: `/workspaces/switch/${reference}`,
        method: "POST",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags updates
     * @name UpdatesList
     * @summary List updates
     * @request GET:/workspaces/updates
     */
    updatesList: (
      query?: {
        /** Page to query data from. Defaults to 1 */
        page?: number;
        /** Number to items to return. Defaults to 10 items */
        per_page?: number;
        /** filter results by the status of the update. */
        status?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ServerListUpdateResponse, ServerAPIStatus>({
        path: `/workspaces/updates`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags updates
     * @name UpdatesCreate
     * @summary Create a new update
     * @request POST:/workspaces/updates
     */
    updatesCreate: (data: ServerCreateUpdateContent, params: RequestParams = {}) =>
      this.request<ServerCreatedUpdateResponse, ServerAPIStatus>({
        path: `/workspaces/updates`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags updates
     * @name DeleteUpdate
     * @summary Delete a specific update
     * @request DELETE:/workspaces/updates/{reference}
     */
    deleteUpdate: (reference: string, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/workspaces/updates/${reference}`,
        method: "DELETE",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags updates
     * @name FetchUpdate
     * @summary Fetch a specific update
     * @request GET:/workspaces/updates/{reference}
     */
    fetchUpdate: (reference: string, params: RequestParams = {}) =>
      this.request<ServerFetchUpdateReponse, ServerAPIStatus>({
        path: `/workspaces/updates/${reference}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags updates
     * @name SendUpdate
     * @summary Send an update to real users
     * @request POST:/workspaces/updates/{reference}
     */
    sendUpdate: (reference: string, data: ServerSendUpdateRequest, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/workspaces/updates/${reference}`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags updates
     * @name UpdateContent
     * @summary Update a specific update
     * @request PUT:/workspaces/updates/{reference}
     */
    updateContent: (reference: string, data: ServerContentUpdateRequest, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/workspaces/updates/${reference}`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags updates
     * @name FetchUpdateAnalytics
     * @summary Fetch analytics for a specific update
     * @request GET:/workspaces/updates/{reference}/analytics
     */
    fetchUpdateAnalytics: (reference: string, params: RequestParams = {}) =>
      this.request<ServerFetchUpdateAnalyticsResponse, ServerAPIStatus>({
        path: `/workspaces/updates/${reference}/analytics`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags updates
     * @name DuplicateUpdate
     * @summary Duplicate a specific update
     * @request POST:/workspaces/updates/{reference}/duplicate
     */
    duplicateUpdate: (reference: string, params: RequestParams = {}) =>
      this.request<ServerCreatedUpdateResponse, ServerAPIStatus>({
        path: `/workspaces/updates/${reference}/duplicate`,
        method: "POST",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags updates
     * @name ToggleUpdatePin
     * @summary Toggle pinned status a specific update
     * @request POST:/workspaces/updates/{reference}/pin
     */
    toggleUpdatePin: (reference: string, params: RequestParams = {}) =>
      this.request<ServerCreatedUpdateResponse, ServerAPIStatus>({
        path: `/workspaces/updates/${reference}/pin`,
        method: "POST",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags updates
     * @name PreviewUpdate
     * @summary Send preview of an update
     * @request POST:/workspaces/updates/{reference}/preview
     */
    previewUpdate: (reference: string, data: ServerPreviewUpdateRequest, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/workspaces/updates/${reference}/preview`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags updates
     * @name UpdatesPinsList
     * @summary List pinned updates
     * @request GET:/workspaces/updates/pins
     */
    updatesPinsList: (params: RequestParams = {}) =>
      this.request<ServerListUpdateResponse, ServerAPIStatus>({
        path: `/workspaces/updates/pins`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
}
