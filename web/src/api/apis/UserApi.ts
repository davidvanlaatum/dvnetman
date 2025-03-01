// noinspection all
// @ts-nocheck
/* tslint:disable */
/* eslint-disable */
/**
 * DVNetMan
 * DVNetMan
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import * as runtime from '../runtime';
import type {
  CurrentUser,
  User,
  UserProvider,
  UserSearchBody,
  UserSearchResults,
} from '../models/index';
import {
    CurrentUserFromJSON,
    CurrentUserToJSON,
    UserFromJSON,
    UserToJSON,
    UserProviderFromJSON,
    UserProviderToJSON,
    UserSearchBodyFromJSON,
    UserSearchBodyToJSON,
    UserSearchResultsFromJSON,
    UserSearchResultsToJSON,
} from '../models/index';

export interface CreateUserRequest {
    user: Omit<User, 'created'|'id'|'updated'>;
}

export interface DeleteUserRequest {
    id: string;
}

export interface GetUserRequest {
    id: string;
    ifNoneMatch?: string;
    ifModifiedSince?: Date;
}

export interface ListUsersRequest {
    page?: number;
    perPage?: number;
    sort?: string;
    userSearchBody?: UserSearchBody;
}

export interface UpdateUserRequest {
    id: string;
    user: Omit<User, 'created'|'id'|'updated'>;
}

/**
 * UserApi - interface
 * 
 * @export
 * @interface UserApiInterface
 */
export interface UserApiInterface {
    /**
     * 
     * @param {User} user 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof UserApiInterface
     */
    createUserRaw(requestParameters: CreateUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<User>>;

    /**
     */
    createUser(requestParameters: CreateUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<User>;

    /**
     * 
     * @param {string} id 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof UserApiInterface
     */
    deleteUserRaw(requestParameters: DeleteUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>>;

    /**
     */
    deleteUser(requestParameters: DeleteUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void>;

    /**
     * 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof UserApiInterface
     */
    getCurrentUserRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CurrentUser>>;

    /**
     */
    getCurrentUser(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CurrentUser>;

    /**
     * 
     * @param {string} id 
     * @param {string} [ifNoneMatch] 
     * @param {Date} [ifModifiedSince] 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof UserApiInterface
     */
    getUserRaw(requestParameters: GetUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<User>>;

    /**
     */
    getUser(requestParameters: GetUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<User>;

    /**
     * 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof UserApiInterface
     */
    getUserProvidersRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<UserProvider>>>;

    /**
     */
    getUserProviders(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<UserProvider>>;

    /**
     * 
     * @param {number} [page] Page number
     * @param {number} [perPage] Number of items per page
     * @param {string} [sort] Sort order
     * @param {UserSearchBody} [userSearchBody] 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof UserApiInterface
     */
    listUsersRaw(requestParameters: ListUsersRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<UserSearchResults>>;

    /**
     */
    listUsers(requestParameters: ListUsersRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<UserSearchResults>;

    /**
     * 
     * @param {string} id 
     * @param {User} user 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof UserApiInterface
     */
    updateUserRaw(requestParameters: UpdateUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<User>>;

    /**
     */
    updateUser(requestParameters: UpdateUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<User>;

}

/**
 * 
 */
export class UserApi extends runtime.BaseAPI implements UserApiInterface {

    /**
     */
    async createUserRaw(requestParameters: CreateUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<User>> {
        if (requestParameters['user'] == null) {
            throw new runtime.RequiredError(
                'user',
                'Required parameter "user" was null or undefined when calling createUser().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/api/v1/user`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: UserToJSON(requestParameters['user']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => UserFromJSON(jsonValue));
    }

    /**
     */
    async createUser(requestParameters: CreateUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<User> {
        const response = await this.createUserRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async deleteUserRaw(requestParameters: DeleteUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters['id'] == null) {
            throw new runtime.RequiredError(
                'id',
                'Required parameter "id" was null or undefined when calling deleteUser().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/api/v1/user/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters['id']))),
            method: 'DELETE',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     */
    async deleteUser(requestParameters: DeleteUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.deleteUserRaw(requestParameters, initOverrides);
    }

    /**
     */
    async getCurrentUserRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CurrentUser>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/api/v1/user/current`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CurrentUserFromJSON(jsonValue));
    }

    /**
     */
    async getCurrentUser(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CurrentUser> {
        const response = await this.getCurrentUserRaw(initOverrides);
        return await response.value();
    }

    /**
     */
    async getUserRaw(requestParameters: GetUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<User>> {
        if (requestParameters['id'] == null) {
            throw new runtime.RequiredError(
                'id',
                'Required parameter "id" was null or undefined when calling getUser().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (requestParameters['ifNoneMatch'] != null) {
            headerParameters['If-None-Match'] = String(requestParameters['ifNoneMatch']);
        }

        if (requestParameters['ifModifiedSince'] != null) {
            headerParameters['If-Modified-Since'] = String(requestParameters['ifModifiedSince']);
        }

        const response = await this.request({
            path: `/api/v1/user/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters['id']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => UserFromJSON(jsonValue));
    }

    /**
     */
    async getUser(requestParameters: GetUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<User> {
        const response = await this.getUserRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async getUserProvidersRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<UserProvider>>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/api/v1/user/providers`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(UserProviderFromJSON));
    }

    /**
     */
    async getUserProviders(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<UserProvider>> {
        const response = await this.getUserProvidersRaw(initOverrides);
        return await response.value();
    }

    /**
     */
    async listUsersRaw(requestParameters: ListUsersRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<UserSearchResults>> {
        const queryParameters: any = {};

        if (requestParameters['page'] != null) {
            queryParameters['page'] = requestParameters['page'];
        }

        if (requestParameters['perPage'] != null) {
            queryParameters['per_page'] = requestParameters['perPage'];
        }

        if (requestParameters['sort'] != null) {
            queryParameters['sort'] = requestParameters['sort'];
        }

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/api/v1/user/search`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: UserSearchBodyToJSON(requestParameters['userSearchBody']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => UserSearchResultsFromJSON(jsonValue));
    }

    /**
     */
    async listUsers(requestParameters: ListUsersRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<UserSearchResults> {
        const response = await this.listUsersRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async updateUserRaw(requestParameters: UpdateUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<User>> {
        if (requestParameters['id'] == null) {
            throw new runtime.RequiredError(
                'id',
                'Required parameter "id" was null or undefined when calling updateUser().'
            );
        }

        if (requestParameters['user'] == null) {
            throw new runtime.RequiredError(
                'user',
                'Required parameter "user" was null or undefined when calling updateUser().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/api/v1/user/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters['id']))),
            method: 'PUT',
            headers: headerParameters,
            query: queryParameters,
            body: UserToJSON(requestParameters['user']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => UserFromJSON(jsonValue));
    }

    /**
     */
    async updateUser(requestParameters: UpdateUserRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<User> {
        const response = await this.updateUserRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
