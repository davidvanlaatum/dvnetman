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

import { mapValues } from '../runtime';
import type { DeviceTypeResult } from './DeviceTypeResult';
import {
    DeviceTypeResultFromJSON,
    DeviceTypeResultFromJSONTyped,
    DeviceTypeResultToJSON,
    DeviceTypeResultToJSONTyped,
} from './DeviceTypeResult';

/**
 * 
 * @export
 * @interface DeviceTypeSearchResults
 */
export interface DeviceTypeSearchResults {
    /**
     * 
     * @type {number}
     * @memberof DeviceTypeSearchResults
     */
    count: number;
    /**
     * 
     * @type {Array<DeviceTypeResult>}
     * @memberof DeviceTypeSearchResults
     */
    items: Array<DeviceTypeResult>;
    /**
     * 
     * @type {boolean}
     * @memberof DeviceTypeSearchResults
     */
    next: boolean;
}

/**
 * Check if a given object implements the DeviceTypeSearchResults interface.
 */
export function instanceOfDeviceTypeSearchResults(value: object): value is DeviceTypeSearchResults {
    if (!('count' in value) || value['count'] === undefined) return false;
    if (!('items' in value) || value['items'] === undefined) return false;
    if (!('next' in value) || value['next'] === undefined) return false;
    return true;
}

export function DeviceTypeSearchResultsFromJSON(json: any): DeviceTypeSearchResults {
    return DeviceTypeSearchResultsFromJSONTyped(json, false);
}

export function DeviceTypeSearchResultsFromJSONTyped(json: any, ignoreDiscriminator: boolean): DeviceTypeSearchResults {
    if (json == null) {
        return json;
    }
    return {
        
        'count': json['count'],
        'items': ((json['items'] as Array<any>).map(DeviceTypeResultFromJSON)),
        'next': json['next'],
    };
}

export function DeviceTypeSearchResultsToJSON(json: any): DeviceTypeSearchResults {
    return DeviceTypeSearchResultsToJSONTyped(json, false);
}

export function DeviceTypeSearchResultsToJSONTyped(value?: DeviceTypeSearchResults | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'count': value['count'],
        'items': ((value['items'] as Array<any>).map(DeviceTypeResultToJSON)),
        'next': value['next'],
    };
}

