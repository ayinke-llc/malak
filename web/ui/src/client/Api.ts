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

export interface MalakPlanMetadata {
  team: {
    enabled: boolean;
    size: number;
  };
}

export enum MalakRole {
  RoleAdmin = "admin",
  RoleMember = "member",
  RoleBilling = "billing",
  RoleInvestor = "investor",
  RoleGuest = "guest",
}

export interface MalakUser {
  created_at: string;
  email: string;
  full_name: string;
  id: string;
  metadata: MalakUserMetadata;
  roles: MalakUserRole[];
  updated_at: string;
}

export interface MalakUserMetadata {
  /**
   * Used to keep track of the last used workspace
   * In the instance of multiple workspaces
   * So when next the user logs in, we remember and take them to the
   * right place rather than always a list of all their workspaces and they
   * have to select one
   */
  current_workspace: string;
}

export interface MalakUserRole {
  created_at: string;
  id: string;
  role: MalakRole;
  updated_at: string;
  user_id: string;
  workspace_id: string;
}

export interface MalakWorkspace {
  created_at: string;
  id: string;
  metadata: MalakPlanMetadata;
  plan_id: string;
  reference: string;
  /**
   * Not required
   * Dummy values work really
   */
  stripe_customer_id: string;
  subscription_id: string;
  updated_at: string;
  workspace_name: string;
}

export interface ServerAPIStatus {
  /** Generic message that tells you the status of the operation */
  message: string;
}

export interface ServerAuthenticateUserRequest {
  code: string;
}

export interface ServerCreateWorkspaceRequest {
  name: string;
}

export interface ServerCreatedUserResponse {
  /** Generic message that tells you the status of the operation */
  message: string;
  token: string;
  user: MalakUser;
  workspace: MalakWorkspace;
}

export interface ServerFetchWorkspaceResponse {
  /** Generic message that tells you the status of the operation */
  message: string;
  workspace: MalakWorkspace;
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
    connectCreate: (provider: string, message: ServerAuthenticateUserRequest, params: RequestParams = {}) =>
      this.request<ServerCreatedUserResponse, ServerAPIStatus>({
        path: `/auth/connect/${provider}`,
        method: "POST",
        body: message,
        type: ContentType.Json,
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
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  workspaces = {
    /**
     * No description
     *
     * @tags workspace
     * @name WorkspacesCreate
     * @summary Create a new workspace
     * @request POST:/workspaces
     */
    workspacesCreate: (message: ServerCreateWorkspaceRequest, params: RequestParams = {}) =>
      this.request<ServerFetchWorkspaceResponse, ServerAPIStatus>({
        path: `/workspaces`,
        method: "POST",
        body: message,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
}
