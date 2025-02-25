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


import * as runtime from '../runtime'
import type { Manufacturer, ManufacturerSearchBody, ManufacturerSearchResults } from '../models/index'
import {
  ManufacturerFromJSON,
  ManufacturerSearchBodyToJSON,
  ManufacturerSearchResultsFromJSON,
  ManufacturerToJSON,
} from '../models/index'

export interface CreateManufacturerRequest {
  manufacturer: Omit<Manufacturer, 'created' | 'id' | 'updated'>
}

export interface DeleteManufacturerRequest {
    id: string;
}

export interface GetManufacturerRequest {
    id: string;
    ifNoneMatch?: string;
    ifModifiedSince?: Date;
}

export interface ListManufacturersRequest {
  page?: number
  perPage?: number
  sort?: string
  manufacturerSearchBody?: ManufacturerSearchBody
}

export interface UpdateManufacturerRequest {
  id: string
  manufacturer: Omit<Manufacturer, 'created' | 'id' | 'updated'>
}

/**
 * ManufacturerApi - interface
 *
 * @export
 * @interface ManufacturerApiInterface
 */
export interface ManufacturerApiInterface {
  /**
   *
   * @param {Manufacturer} manufacturer
   * @param {*} [options] Override http request option.
   * @throws {RequiredError}
   * @memberof ManufacturerApiInterface
   */
  createManufacturerRaw(
    requestParameters: CreateManufacturerRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<Manufacturer>>

  /**
   */
  createManufacturer(
    requestParameters: CreateManufacturerRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<Manufacturer>

  /**
   *
   * @param {string} id
   * @param {*} [options] Override http request option.
   * @throws {RequiredError}
   * @memberof ManufacturerApiInterface
   */
  deleteManufacturerRaw(
    requestParameters: DeleteManufacturerRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<void>>

  /**
   */
  deleteManufacturer(
    requestParameters: DeleteManufacturerRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<void>

  /**
   *
   * @param {string} id
   * @param {string} [ifNoneMatch]
   * @param {Date} [ifModifiedSince]
   * @param {*} [options] Override http request option.
   * @throws {RequiredError}
   * @memberof ManufacturerApiInterface
   */
  getManufacturerRaw(
    requestParameters: GetManufacturerRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<Manufacturer>>

  /**
   */
  getManufacturer(
    requestParameters: GetManufacturerRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<Manufacturer>

  /**
   *
   * @param {number} [page] Page number
   * @param {number} [perPage] Number of items per page
   * @param {string} [sort] Sort order
   * @param {ManufacturerSearchBody} [manufacturerSearchBody]
   * @param {*} [options] Override http request option.
   * @throws {RequiredError}
   * @memberof ManufacturerApiInterface
   */
  listManufacturersRaw(
    requestParameters: ListManufacturersRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<ManufacturerSearchResults>>

  /**
   */
  listManufacturers(
    requestParameters: ListManufacturersRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<ManufacturerSearchResults>

  /**
   *
   * @param {string} id
   * @param {Manufacturer} manufacturer
   * @param {*} [options] Override http request option.
   * @throws {RequiredError}
   * @memberof ManufacturerApiInterface
   */
  updateManufacturerRaw(
    requestParameters: UpdateManufacturerRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<Manufacturer>>

  /**
   */
  updateManufacturer(
    requestParameters: UpdateManufacturerRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<Manufacturer>
}

/**
 *
 */
export class ManufacturerApi extends runtime.BaseAPI implements ManufacturerApiInterface {

    /**
     */
    async createManufacturerRaw(requestParameters: CreateManufacturerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Manufacturer>> {
        if (requestParameters['manufacturer'] == null) {
            throw new runtime.RequiredError(
                'manufacturer',
                'Required parameter "manufacturer" was null or undefined when calling createManufacturer().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/api/v1/manufacturer`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: ManufacturerToJSON(requestParameters['manufacturer']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ManufacturerFromJSON(jsonValue));
    }

    /**
     */
    async createManufacturer(requestParameters: CreateManufacturerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Manufacturer> {
        const response = await this.createManufacturerRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async deleteManufacturerRaw(requestParameters: DeleteManufacturerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters['id'] == null) {
            throw new runtime.RequiredError(
                'id',
                'Required parameter "id" was null or undefined when calling deleteManufacturer().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/api/v1/manufacturer/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters['id']))),
            method: 'DELETE',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     */
    async deleteManufacturer(requestParameters: DeleteManufacturerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.deleteManufacturerRaw(requestParameters, initOverrides);
    }

    /**
     */
    async getManufacturerRaw(requestParameters: GetManufacturerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Manufacturer>> {
        if (requestParameters['id'] == null) {
            throw new runtime.RequiredError(
                'id',
                'Required parameter "id" was null or undefined when calling getManufacturer().'
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
            path: `/api/v1/manufacturer/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters['id']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ManufacturerFromJSON(jsonValue));
    }

    /**
     */
    async getManufacturer(requestParameters: GetManufacturerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Manufacturer> {
        const response = await this.getManufacturerRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async listManufacturersRaw(requestParameters: ListManufacturersRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ManufacturerSearchResults>> {
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

      headerParameters['Content-Type'] = 'application/json'

      constresponse = await this.request(
        {
          path: `/api/v1/manufacturer/search`,
          method: 'POST',
          headers: headerParameters,
          query: queryParameters,
          body: ManufacturerSearchBodyToJSON(requestParameters['manufacturerSearchBody']),
        },
        initOverrides,
      )

        return new runtime.JSONApiResponse(response, (jsonValue) => ManufacturerSearchResultsFromJSON(jsonValue));
    }

    /**
     */
    async listManufacturers(requestParameters: ListManufacturersRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ManufacturerSearchResults> {
        const response = await this.listManufacturersRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     */
    async updateManufacturerRaw(requestParameters: UpdateManufacturerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Manufacturer>> {
        if (requestParameters['id'] == null) {
            throw new runtime.RequiredError(
                'id',
                'Required parameter "id" was null or undefined when calling updateManufacturer().'
            );
        }

        if (requestParameters['manufacturer'] == null) {
            throw new runtime.RequiredError(
                'manufacturer',
                'Required parameter "manufacturer" was null or undefined when calling updateManufacturer().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/api/v1/manufacturer/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters['id']))),
            method: 'PUT',
            headers: headerParameters,
            query: queryParameters,
            body: ManufacturerToJSON(requestParameters['manufacturer']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ManufacturerFromJSON(jsonValue));
    }

    /**
     */
    async updateManufacturer(requestParameters: UpdateManufacturerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Manufacturer> {
        const response = await this.updateManufacturerRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
