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
  Device,
  DeviceSearchBody,
  DeviceSearchResults,
} from '../models/index';
import {
    DeviceFromJSON,
    DeviceToJSON,
    DeviceSearchBodyFromJSON,
    DeviceSearchBodyToJSON,
    DeviceSearchResultsFromJSON,
    DeviceSearchResultsToJSON,
} from '../models/index';

export interface CreateDeviceRequest {
    device: Omit<Device, 'created'|'id'|'updated'>;
}

export interface DeleteDeviceRequest {
    id: string;
}

export interface GetDeviceRequest {
    id: string;
    ifNoneMatch?: string;
    ifModifiedSince?: Date;
}

export interface ListDevicesRequest {
    page?: number;
    perPage?: number;
    sort?: string;
    deviceSearchBody?: DeviceSearchBody;
}

export interface UpdateDeviceRequest {
    id: string;
    device: Omit<Device, 'created'|'id'|'updated'>;
}

/**
 * DeviceApi - interface
 * 
 * @export
 * @interface DeviceApiInterface
 */
export interface DeviceApiInterface {
    /**
     * 
     * @param {Device} device 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof DeviceApiInterface
     */
    createDeviceRaw(requestParameters: CreateDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Device>>;

    /**
     */
    createDevice(requestParameters: CreateDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Device>;

    /**
     * 
     * @param {string} id 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof DeviceApiInterface
     */
    deleteDeviceRaw(requestParameters: DeleteDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>>;

    /**
     */
    deleteDevice(requestParameters: DeleteDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void>;

    /**
     * 
     * @param {string} id 
     * @param {string} [ifNoneMatch] 
     * @param {Date} [ifModifiedSince] 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof DeviceApiInterface
     */
    getDeviceRaw(requestParameters: GetDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Device>>;

    /**
     */
    getDevice(requestParameters: GetDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Device>;

    /**
     * 
     * @param {number} [page] Page number
     * @param {number} [perPage] Number of items per page
     * @param {string} [sort] Sort order
     * @param {DeviceSearchBody} [deviceSearchBody] 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof DeviceApiInterface
     */
    listDevicesRaw(requestParameters: ListDevicesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<DeviceSearchResults>>;

    /**
     */
    listDevices(requestParameters: ListDevicesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<DeviceSearchResults>;

    /**
     * 
     * @param {string} id 
     * @param {Device} device 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof DeviceApiInterface
     */
    updateDeviceRaw(requestParameters: UpdateDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Device>>;

    /**
     */
    updateDevice(requestParameters: UpdateDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Device>;

}

/**
 * 
 */
export class DeviceApi extends runtime.BaseAPI implements DeviceApiInterface {

    /**
     */
    async createDeviceRaw(requestParameters: CreateDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Device>> {
        if (requestParameters['device'] == null) {
            throw new runtime.RequiredError(
                'device',
                'Required parameter "device" was null or undefined when calling createDevice().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/api/v1/device`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: DeviceToJSON(requestParameters['device']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => DeviceFromJSON(jsonValue));
    }

    /**
     */
    async createDevice(requestParameters: CreateDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Device> {
        const response = await this.createDeviceRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async deleteDeviceRaw(requestParameters: DeleteDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters['id'] == null) {
            throw new runtime.RequiredError(
                'id',
                'Required parameter "id" was null or undefined when calling deleteDevice().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/api/v1/device/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters['id']))),
            method: 'DELETE',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     */
    async deleteDevice(requestParameters: DeleteDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.deleteDeviceRaw(requestParameters, initOverrides);
    }

    /**
     */
    async getDeviceRaw(requestParameters: GetDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Device>> {
        if (requestParameters['id'] == null) {
            throw new runtime.RequiredError(
                'id',
                'Required parameter "id" was null or undefined when calling getDevice().'
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
            path: `/api/v1/device/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters['id']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => DeviceFromJSON(jsonValue));
    }

    /**
     */
    async getDevice(requestParameters: GetDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Device> {
        const response = await this.getDeviceRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async listDevicesRaw(requestParameters: ListDevicesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<DeviceSearchResults>> {
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
            path: `/api/v1/device/search`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: DeviceSearchBodyToJSON(requestParameters['deviceSearchBody']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => DeviceSearchResultsFromJSON(jsonValue));
    }

    /**
     */
    async listDevices(requestParameters: ListDevicesRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<DeviceSearchResults> {
        const response = await this.listDevicesRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async updateDeviceRaw(requestParameters: UpdateDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Device>> {
        if (requestParameters['id'] == null) {
            throw new runtime.RequiredError(
                'id',
                'Required parameter "id" was null or undefined when calling updateDevice().'
            );
        }

        if (requestParameters['device'] == null) {
            throw new runtime.RequiredError(
                'device',
                'Required parameter "device" was null or undefined when calling updateDevice().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/api/v1/device/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters['id']))),
            method: 'PUT',
            headers: headerParameters,
            query: queryParameters,
            body: DeviceToJSON(requestParameters['device']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => DeviceFromJSON(jsonValue));
    }

    /**
     */
    async updateDevice(requestParameters: UpdateDeviceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Device> {
        const response = await this.updateDeviceRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
