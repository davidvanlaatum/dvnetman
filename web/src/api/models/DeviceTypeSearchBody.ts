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
/**
 * 
 * @export
 * @interface DeviceTypeSearchBody
 */
export interface DeviceTypeSearchBody {
    /**
     * 
     * @type {Array<string>}
     * @memberof DeviceTypeSearchBody
     */
    fields?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof DeviceTypeSearchBody
     */
    ids?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof DeviceTypeSearchBody
     */
    manufacturer?: Array<string>;
    /**
     * 
     * @type {string}
     * @memberof DeviceTypeSearchBody
     */
    model?: string;
    /**
     * 
     * @type {string}
     * @memberof DeviceTypeSearchBody
     */
    modelRegex?: string;
}

/**
 * Check if a given object implements the DeviceTypeSearchBody interface.
 */
export function instanceOfDeviceTypeSearchBody(value: object): value is DeviceTypeSearchBody {
    return true;
}

export function DeviceTypeSearchBodyFromJSON(json: any): DeviceTypeSearchBody {
    return DeviceTypeSearchBodyFromJSONTyped(json, false);
}

export function DeviceTypeSearchBodyFromJSONTyped(json: any, ignoreDiscriminator: boolean): DeviceTypeSearchBody {
    if (json == null) {
        return json;
    }
    return {
        
        'fields': json['fields'] == null ? undefined : json['fields'],
        'ids': json['ids'] == null ? undefined : json['ids'],
        'manufacturer': json['manufacturer'] == null ? undefined : json['manufacturer'],
        'model': json['model'] == null ? undefined : json['model'],
        'modelRegex': json['modelRegex'] == null ? undefined : json['modelRegex'],
    };
}

export function DeviceTypeSearchBodyToJSON(json: any): DeviceTypeSearchBody {
    return DeviceTypeSearchBodyToJSONTyped(json, false);
}

export function DeviceTypeSearchBodyToJSONTyped(value?: DeviceTypeSearchBody | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'fields': value['fields'],
        'ids': value['ids'],
        'manufacturer': value['manufacturer'],
        'model': value['model'],
        'modelRegex': value['modelRegex'],
    };
}

