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

export interface MalakAPIKey {
  created_at?: string;
  created_by?: string;
  expires_at?: string;
  id?: string;
  key_name?: string;
  reference?: string;
  updated_at?: string;
  workspace_id?: string;
}

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

export interface MalakDashboardChart {
  chart?: MalakIntegrationChart;
  chart_id?: string;
  created_at?: string;
  dashboard_id?: string;
  id?: string;
  reference?: string;
  updated_at?: string;
  workspace_id?: string;
  workspace_integration_id?: string;
}

export interface MalakDashboardChartPosition {
  chart_id?: string;
  dashboard_id?: string;
  id?: string;
  order_index?: number;
}

export interface MalakDashboardLink {
  contact?: MalakContact;
  contact_id?: string;
  created_at?: string;
  dashboard?: MalakDashboard;
  dashboard_id?: string;
  expires_at?: string;
  id?: string;
  link_type?: MalakDashboardLinkType;
  reference?: string;
  token?: string;
  updated_at?: string;
}

export enum MalakDashboardLinkType {
  DashboardLinkTypeDefault = "default",
  DashboardLinkTypeContact = "contact",
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

export interface MalakDeckDailyEngagement {
  created_at?: string;
  deck_id?: string;
  engagement_count?: number;
  engagement_date?: string;
  id?: string;
  reference?: string;
  updated_at?: string;
  workspace_id?: string;
}

export interface MalakDeckEngagementResponse {
  daily_engagements: MalakDeckDailyEngagement[];
  geographic_stats: MalakDeckGeographicStat[];
}

export interface MalakDeckGeographicStat {
  country?: string;
  created_at?: string;
  deck_id?: string;
  id?: string;
  reference?: string;
  stat_date?: string;
  updated_at?: string;
  view_count?: number;
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

export interface MalakDeckViewerSession {
  browser?: string;
  city?: string;
  contact?: MalakContact;
  contact_id?: string;
  country?: string;
  created_at?: string;
  deck_id?: string;
  device_info?: string;
  id?: string;
  ip_address?: string;
  os?: string;
  reference?: string;
  session_id?: string;
  time_spent_seconds?: number;
  updated_at?: string;
  viewed_at?: string;
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
  chart_type?: MalakIntegrationChartType;
  created_at?: string;
  data_point_type?: MalakIntegrationDataPointType;
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

export enum MalakIntegrationChartType {
  IntegrationChartTypeBar = "bar",
  IntegrationChartTypePie = "pie",
}

export interface MalakIntegrationDataPoint {
  created_at?: string;
  id?: string;
  integration_chart_id?: string;
  metadata?: MalakIntegrationDataPointMetadata;
  point_name?: string;
  point_value?: number;
  reference?: string;
  updated_at?: string;
  workspace_id?: string;
  workspace_integration_id?: string;
}

export type MalakIntegrationDataPointMetadata = object;

export enum MalakIntegrationDataPointType {
  IntegrationDataPointTypeCurrency = "currency",
  IntegrationDataPointTypeOthers = "others",
}

export interface MalakIntegrationMetadata {
  endpoint?: string;
}

export enum MalakIntegrationType {
  IntegrationTypeOauth2 = "oauth2",
  IntegrationTypeApiKey = "api_key",
  IntegrationTypeSystem = "system",
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
    max_charts_per_dashboard?: number;
    share_dashboard_via_link?: boolean;
  };
  data_room?: {
    share_via_link?: boolean;
    size?: number;
  };
  deck?: {
    analytics?: {
      can_view_historical_sessions?: boolean;
    };
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
  preferences?: MalakPublicDeckPreference;
  reference?: string;
  session?: MalakDeckViewerSession;
  short_link?: string;
  title?: string;
  updated_at?: string;
  workspace_id?: string;
}

export interface MalakPublicDeckPreference {
  enable_downloading?: boolean;
  has_password?: boolean;
  require_email?: boolean;
}

export enum MalakRecipientStatus {
  RecipientStatusPending = "pending",
  RecipientStatusSent = "sent",
  RecipientStatusFailed = "failed",
}

export enum MalakRevocationType {
  RevocationTypeImmediate = "immediate",
  RevocationTypeDay = "day",
  RevocationTypeWeek = "week",
}

export enum MalakRole {
  RoleAdmin = "admin",
  RoleMember = "member",
  RoleBilling = "billing",
  RoleInvestor = "investor",
  RoleGuest = "guest",
}

export interface MalakSystemTemplate {
  content?: MalakBlock[];
  created_at?: string;
  description?: string;
  id?: string;
  number_of_uses?: number;
  reference?: string;
  title?: string;
  updated_at?: string;
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

export interface ServerAddChartToDashboardRequest {
  chart_reference: string;
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

export interface ServerCreateAPIKeyRequest {
  title: string;
}

export interface ServerCreateChartRequest {
  chart_type: MalakIntegrationChartType;
  title: string;
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
  description: string;
  title: string;
}

export interface ServerCreateDeckRequest {
  deck_url?: string;
  title?: string;
}

export interface ServerCreateDeckViewerSession {
  browser: string;
  device_info: string;
  os: string;
  password: string;
}

export interface ServerCreateUpdateContent {
  template?: {
    is_system_template?: boolean;
    reference?: string;
  };
  title: string;
}

export interface ServerCreateWorkspaceRequest {
  name: string;
}

export interface ServerCreatedAPIKeyResponse {
  message: string;
  value: string;
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

export interface ServerFetchEngagementsResponse {
  engagements: MalakDeckEngagementResponse;
  message: string;
}

export interface ServerFetchPublicDeckResponse {
  deck: MalakPublicDeck;
  message: string;
}

export interface ServerFetchSessionsDeck {
  message: string;
  meta: ServerMeta;
  sessions: MalakDeckViewerSession[];
}

export interface ServerFetchTemplatesResponse {
  message: string;
  templates: {
    system: MalakSystemTemplate[];
    workspace: MalakSystemTemplate[];
  };
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

export interface ServerGenerateDashboardLinkRequest {
  email?: string;
}

export interface ServerListAPIKeysResponse {
  keys: MalakAPIKey[];
  message: string;
}

export interface ServerListChartDataPointsResponse {
  data_points: MalakIntegrationDataPoint[];
  message: string;
}

export interface ServerListContactsResponse {
  contacts: MalakContact[];
  message: string;
  meta: ServerMeta;
}

export interface ServerListDashboardChartsResponse {
  charts: MalakDashboardChart[];
  dashboard: MalakDashboard;
  link: MalakDashboardLink;
  message: string;
  positions: MalakDashboardChartPosition[];
}

export interface ServerListDashboardLinkResponse {
  links: MalakDashboardLink[];
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

export interface ServerRegenerateLinkResponse {
  link: MalakDashboardLink;
  message: string;
}

export interface ServerRevokeAPIKeyRequest {
  strategy: MalakRevocationType;
}

export interface ServerSendUpdateRequest {
  emails?: string[];
  send_at?: number;
}

export interface ServerTestAPIIntegrationRequest {
  api_key: string;
}

export interface ServerUpdateDashboardPositionsRequest {
  positions: {
    chart_id: string;
    index: number;
  }[];
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
    this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "" });
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
 * @contact Ayinke Ventures <lanre@ayinke.ventures>
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  auth = {
    /**
     * @description Sign in with a social login provider
     *
     * @tags auth
     * @name ConnectCreate
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
     * @description list your contacts
     *
     * @tags contacts
     * @name ContactsList
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
     * @description Creates a new contact
     *
     * @tags contacts
     * @name ContactsCreate
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
     * @description delete a contact
     *
     * @tags contacts
     * @name DeleteContact
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
     * @description fetch a contact by reference
     *
     * @tags contacts
     * @name ContactsDetail
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
     * @description edit a contact
     *
     * @tags contacts
     * @name ContactsUpdate
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
     * @description List all created contact lists
     *
     * @tags contacts
     * @name FetchContactLists
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
     * @description Create a new contact list
     *
     * @tags contacts
     * @name CreateContactList
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
     * @description delete a contact list
     *
     * @tags contacts
     * @name DeleteContactList
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
     * @description add a new contact to a list
     *
     * @tags contacts
     * @name AddEmailToContactList
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
     * @description Edit a contact list
     *
     * @tags contacts
     * @name EditContactList
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
     * @description List dashboards
     *
     * @tags dashboards
     * @name DashboardsList
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
     * @description create a new dashboard
     *
     * @tags dashboards
     * @name DashboardsCreate
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
     * @description fetch dashboard
     *
     * @tags dashboards
     * @name DashboardsDetail
     * @request GET:/dashboards/{reference}
     */
    dashboardsDetail: (reference: string, params: RequestParams = {}) =>
      this.request<ServerListDashboardChartsResponse, ServerAPIStatus>({
        path: `/dashboards/${reference}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description list access controls
     *
     * @tags dashboards
     * @name AccessControlDetail
     * @request GET:/dashboards/{reference}/access-control
     */
    accessControlDetail: (
      reference: string,
      query?: {
        /** Page to query data from. Defaults to 1 */
        page?: number;
        /** Number to items to return. Defaults to 10 items */
        per_page?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<ServerListDashboardLinkResponse, ServerAPIStatus>({
        path: `/dashboards/${reference}/access-control`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * @description delete access controls
     *
     * @tags dashboards
     * @name AccessControlDelete
     * @request DELETE:/dashboards/{reference}/access-control/{link_reference}
     */
    accessControlDelete: (reference: string, linkReference: string, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/dashboards/${reference}/access-control/${linkReference}`,
        method: "DELETE",
        format: "json",
        ...params,
      }),

    /**
     * @description regenerate the default link for a dashboard
     *
     * @tags dashboards
     * @name AccessControlLinkCreate
     * @request POST:/dashboards/{reference}/access-control/link
     */
    accessControlLinkCreate: (
      reference: string,
      data: ServerGenerateDashboardLinkRequest,
      params: RequestParams = {},
    ) =>
      this.request<ServerRegenerateLinkResponse, ServerAPIStatus>({
        path: `/dashboards/${reference}/access-control/link`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description remove a chart from a dashboard
     *
     * @tags dashboards
     * @name ChartsDelete
     * @request DELETE:/dashboards/{reference}/charts
     */
    chartsDelete: (reference: string, data: ServerAddChartToDashboardRequest, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/dashboards/${reference}/charts`,
        method: "DELETE",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description add a chart to a dashboard
     *
     * @tags dashboards
     * @name ChartsUpdate
     * @request PUT:/dashboards/{reference}/charts
     */
    chartsUpdate: (reference: string, data: ServerAddChartToDashboardRequest, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/dashboards/${reference}/charts`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description update dashboard chart positioning
     *
     * @tags dashboards
     * @name PositionsCreate
     * @request POST:/dashboards/{reference}/positions
     */
    positionsCreate: (reference: string, data: ServerUpdateDashboardPositionsRequest, params: RequestParams = {}) =>
      this.request<any, ServerAPIStatus>({
        path: `/dashboards/${reference}/positions`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description List charts
     *
     * @tags dashboards
     * @name ChartsList
     * @request GET:/dashboards/charts
     */
    chartsList: (params: RequestParams = {}) =>
      this.request<ServerListIntegrationChartsResponse, ServerAPIStatus>({
        path: `/dashboards/charts`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description fetch charting data
     *
     * @tags dashboards
     * @name ChartsDetail
     * @request GET:/dashboards/charts/{reference}
     */
    chartsDetail: (reference: string, params: RequestParams = {}) =>
      this.request<ServerListChartDataPointsResponse, ServerAPIStatus>({
        path: `/dashboards/charts/${reference}`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
  decks = {
    /**
     * @description list all decks. No pagination
     *
     * @tags decks
     * @name DecksList
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
     * @description Creates a new deck
     *
     * @tags decks
     * @name DecksCreate
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
     * @description delete a deck
     *
     * @tags decks
     * @name DecksDelete
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
     * @description fetch a deck
     *
     * @tags decks
     * @name DecksDetail
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
     * @description fetch deck engagements and geographic stats
     *
     * @tags decks
     * @name AnalyticsDetail
     * @request GET:/decks/{reference}/analytics
     */
    analyticsDetail: (reference: string, params: RequestParams = {}) =>
      this.request<ServerFetchEngagementsResponse, ServerAPIStatus>({
        path: `/decks/${reference}/analytics`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description toggle archive status of a deck
     *
     * @tags decks
     * @name ToggleArchive
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
     * @description toggle pinned status of a deck
     *
     * @tags decks
     * @name TogglePin
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
     * @description update a deck preferences
     *
     * @tags decks
     * @name PreferencesUpdate
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

    /**
     * @description fetch deck viewing sessions on dashboard
     *
     * @tags decks
     * @name SessionsDetail
     * @request GET:/decks/{reference}/sessions
     */
    sessionsDetail: (
      reference: string,
      query?: {
        /** Page to query data from. Defaults to 1 */
        page?: number;
        /** Number to items to return. Defaults to 10 items */
        per_page?: number;
        /** number of days to fetch deck sessions */
        days?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<ServerFetchSessionsDeck, ServerAPIStatus>({
        path: `/decks/${reference}/sessions`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),
  };
  developers = {
    /**
     * @description list api keys
     *
     * @tags developers
     * @name KeysList
     * @request GET:/developers/keys
     */
    keysList: (params: RequestParams = {}) =>
      this.request<ServerListAPIKeysResponse, ServerAPIStatus>({
        path: `/developers/keys`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description Creates a new api key
     *
     * @tags developers
     * @name KeysCreate
     * @request POST:/developers/keys
     */
    keysCreate: (data: ServerCreateAPIKeyRequest, params: RequestParams = {}) =>
      this.request<ServerCreatedAPIKeyResponse, ServerAPIStatus>({
        path: `/developers/keys`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description revoke a specific api key
     *
     * @tags developers
     * @name KeysDelete
     * @request DELETE:/developers/keys/{reference}
     */
    keysDelete: (reference: string, data: ServerRevokeAPIKeyRequest, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/developers/keys/${reference}`,
        method: "DELETE",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  public = {
    /**
     * @description fetch public dashboard and charting data points
     *
     * @tags dashboards
     * @name DashboardsDetail
     * @request GET:/public/dashboards/{reference}
     */
    dashboardsDetail: (reference: string, params: RequestParams = {}) =>
      this.request<ServerListDashboardChartsResponse, ServerAPIStatus>({
        path: `/public/dashboards/${reference}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description fetch charting data
     *
     * @tags dashboards
     * @name DashboardsChartsDetail
     * @request GET:/public/dashboards/{reference}/charts/{chart_reference}
     */
    dashboardsChartsDetail: (reference: string, chartReference: string, params: RequestParams = {}) =>
      this.request<ServerListChartDataPointsResponse, ServerAPIStatus>({
        path: `/public/dashboards/${reference}/charts/${chartReference}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description public api to fetch a deck
     *
     * @tags decks-viewer
     * @name DecksCreate
     * @request POST:/public/decks/{reference}
     */
    decksCreate: (reference: string, data: ServerCreateDeckViewerSession, params: RequestParams = {}) =>
      this.request<ServerFetchPublicDeckResponse, ServerAPIStatus>({
        path: `/public/decks/${reference}`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description update the session details
     *
     * @tags decks-viewer
     * @name DecksUpdate
     * @request PUT:/public/decks/{reference}
     */
    decksUpdate: (reference: string, data: ServerCreateDeckViewerSession, params: RequestParams = {}) =>
      this.request<ServerFetchPublicDeckResponse, ServerAPIStatus>({
        path: `/public/decks/${reference}`,
        method: "PUT",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  updates = {
    /**
     * @description Fetch a specific update
     *
     * @tags updates
     * @name ReactPost
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
     * @description Upload a deck
     *
     * @tags decks
     * @name UploadDeck
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
     * @description Upload an image
     *
     * @tags images
     * @name UploadImage
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
     * @description Fetch current user. This api should also double as a token validation api
     *
     * @tags user
     * @name UserList
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
     * @description update workspace details
     *
     * @tags workspace
     * @name WorkspacesPartialUpdate
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
     * @description Create a new workspace
     *
     * @tags workspace
     * @name WorkspacesCreate
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
     * @description get billing portal
     *
     * @tags billing
     * @name BillingCreate
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
     * @description fetch workspace preferences
     *
     * @tags integrations
     * @name IntegrationsList
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
     * @description disable integration
     *
     * @tags integrations
     * @name IntegrationsDelete
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
     * @description enable integration
     *
     * @tags integrations
     * @name IntegrationsCreate
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
     * @description update integration api key
     *
     * @tags integrations
     * @name IntegrationsUpdate
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
     * @description create chart
     *
     * @tags integrations
     * @name IntegrationsChartsCreate
     * @request POST:/workspaces/integrations/{reference}/charts
     */
    integrationsChartsCreate: (reference: string, data: ServerCreateChartRequest, params: RequestParams = {}) =>
      this.request<ServerAPIStatus, ServerAPIStatus>({
        path: `/workspaces/integrations/${reference}/charts`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description test an api key is valid and can reach the integration
     *
     * @tags integrations
     * @name IntegrationsPingCreate
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
     * @description fetch workspace preferences
     *
     * @tags workspace
     * @name PreferencesList
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
     * @description update workspace preferences
     *
     * @tags workspace
     * @name PreferencesUpdate
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
     * @description Switch current workspace
     *
     * @tags workspace
     * @name Switchworkspace
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
     * @description List updates
     *
     * @tags updates
     * @name UpdatesList
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
     * @description Create a new update
     *
     * @tags updates
     * @name UpdatesCreate
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
     * @description Delete a specific update
     *
     * @tags updates
     * @name DeleteUpdate
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
     * @description Fetch a specific update
     *
     * @tags updates
     * @name FetchUpdate
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
     * @description Send an update to real users
     *
     * @tags updates
     * @name SendUpdate
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
     * @description Update a specific update
     *
     * @tags updates
     * @name UpdateContent
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
     * @description Fetch analytics for a specific update
     *
     * @tags updates
     * @name FetchUpdateAnalytics
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
     * @description Duplicate a specific update
     *
     * @tags updates
     * @name DuplicateUpdate
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
     * @description Toggle pinned status a specific update
     *
     * @tags updates
     * @name ToggleUpdatePin
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
     * @description Send preview of an update
     *
     * @tags updates
     * @name PreviewUpdate
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
     * @description List pinned updates
     *
     * @tags updates
     * @name UpdatesPinsList
     * @request GET:/workspaces/updates/pins
     */
    updatesPinsList: (params: RequestParams = {}) =>
      this.request<ServerListUpdateResponse, ServerAPIStatus>({
        path: `/workspaces/updates/pins`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description list all templates. this will include both systems and your own created templates
     *
     * @tags updates
     * @name UpdatesTemplatesList
     * @request GET:/workspaces/updates/templates
     */
    updatesTemplatesList: (params: RequestParams = {}) =>
      this.request<ServerFetchTemplatesResponse, ServerAPIStatus>({
        path: `/workspaces/updates/templates`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
}
