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
  Location,
  LocationSearchBody,
  LocationSearchResults,
} from '../models/index';
import {
    LocationFromJSON,
    LocationToJSON,
    LocationSearchBodyFromJSON,
    LocationSearchBodyToJSON,
    LocationSearchResultsFromJSON,
    LocationSearchResultsToJSON,
} from '../models/index';

export interface CreateLocationRequest {
    location: Omit<Location, 'created'|'id'|'updated'>;
}

export interface DeleteLocationRequest {
    id: string;
}

export interface GetLocationRequest {
    id: string;
    ifNoneMatch?: string;
    ifModifiedSince?: Date;
}

export interface ListLocationsRequest {
    page?: number;
    perPage?: number;
    sort?: string;
    locationSearchBody?: LocationSearchBody;
}

export interface UpdateLocationRequest {
    id: string;
    location: Omit<Location, 'created'|'id'|'updated'>;
}

/**
 * LocationApi - interface
 * 
 * @export
 * @interface LocationApiInterface
 */
export interface LocationApiInterface {
    /**
     * 
     * @param {Location} location 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof LocationApiInterface
     */
    createLocationRaw(requestParameters: CreateLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Location>>;

    /**
     */
    createLocation(requestParameters: CreateLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Location>;

    /**
     * 
     * @param {string} id 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof LocationApiInterface
     */
    deleteLocationRaw(requestParameters: DeleteLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>>;

    /**
     */
    deleteLocation(requestParameters: DeleteLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void>;

    /**
     * 
     * @param {string} id 
     * @param {string} [ifNoneMatch] 
     * @param {Date} [ifModifiedSince] 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof LocationApiInterface
     */
    getLocationRaw(requestParameters: GetLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Location>>;

    /**
     */
    getLocation(requestParameters: GetLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Location>;

    /**
     * 
     * @param {number} [page] Page number
     * @param {number} [perPage] Number of items per page
     * @param {string} [sort] Sort order
     * @param {LocationSearchBody} [locationSearchBody] 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof LocationApiInterface
     */
    listLocationsRaw(requestParameters: ListLocationsRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<LocationSearchResults>>;

    /**
     */
    listLocations(requestParameters: ListLocationsRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<LocationSearchResults>;

    /**
     * 
     * @param {string} id 
     * @param {Location} location 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof LocationApiInterface
     */
    updateLocationRaw(requestParameters: UpdateLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Location>>;

    /**
     */
    updateLocation(requestParameters: UpdateLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Location>;

}

/**
 * 
 */
export class LocationApi extends runtime.BaseAPI implements LocationApiInterface {

    /**
     */
    async createLocationRaw(requestParameters: CreateLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Location>> {
        if (requestParameters['location'] == null) {
            throw new runtime.RequiredError(
                'location',
                'Required parameter "location" was null or undefined when calling createLocation().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/api/v1/location`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: LocationToJSON(requestParameters['location']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => LocationFromJSON(jsonValue));
    }

    /**
     */
    async createLocation(requestParameters: CreateLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Location> {
        const response = await this.createLocationRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async deleteLocationRaw(requestParameters: DeleteLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters['id'] == null) {
            throw new runtime.RequiredError(
                'id',
                'Required parameter "id" was null or undefined when calling deleteLocation().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/api/v1/location/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters['id']))),
            method: 'DELETE',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     */
    async deleteLocation(requestParameters: DeleteLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.deleteLocationRaw(requestParameters, initOverrides);
    }

    /**
     */
    async getLocationRaw(requestParameters: GetLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Location>> {
        if (requestParameters['id'] == null) {
            throw new runtime.RequiredError(
                'id',
                'Required parameter "id" was null or undefined when calling getLocation().'
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
            path: `/api/v1/location/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters['id']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => LocationFromJSON(jsonValue));
    }

    /**
     */
    async getLocation(requestParameters: GetLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Location> {
        const response = await this.getLocationRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async listLocationsRaw(requestParameters: ListLocationsRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<LocationSearchResults>> {
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
            path: `/api/v1/location/search`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: LocationSearchBodyToJSON(requestParameters['locationSearchBody']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => LocationSearchResultsFromJSON(jsonValue));
    }

    /**
     */
    async listLocations(requestParameters: ListLocationsRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<LocationSearchResults> {
        const response = await this.listLocationsRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async updateLocationRaw(requestParameters: UpdateLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Location>> {
        if (requestParameters['id'] == null) {
            throw new runtime.RequiredError(
                'id',
                'Required parameter "id" was null or undefined when calling updateLocation().'
            );
        }

        if (requestParameters['location'] == null) {
            throw new runtime.RequiredError(
                'location',
                'Required parameter "location" was null or undefined when calling updateLocation().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/api/v1/location/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters['id']))),
            method: 'PUT',
            headers: headerParameters,
            query: queryParameters,
            body: LocationToJSON(requestParameters['location']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => LocationFromJSON(jsonValue));
    }

    /**
     */
    async updateLocation(requestParameters: UpdateLocationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Location> {
        const response = await this.updateLocationRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
