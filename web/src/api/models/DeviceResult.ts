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
import type { ObjectReference } from './ObjectReference';
import {
    ObjectReferenceFromJSON,
    ObjectReferenceFromJSONTyped,
    ObjectReferenceToJSON,
    ObjectReferenceToJSONTyped,
} from './ObjectReference';
import type { Tag } from './Tag';
import {
    TagFromJSON,
    TagFromJSONTyped,
    TagToJSON,
    TagToJSONTyped,
} from './Tag';

/**
 * 
 * @export
 * @interface DeviceResult
 */
export interface DeviceResult {
    /**
     * 
     * @type {Date}
     * @memberof DeviceResult
     */
    created?: Date;
    /**
     * 
     * @type {string}
     * @memberof DeviceResult
     */
    description?: string;
    /**
     * 
     * @type {ObjectReference}
     * @memberof DeviceResult
     */
    deviceType?: ObjectReference;
    /**
     * 
     * @type {string}
     * @memberof DeviceResult
     */
    readonly id: string;
    /**
     * 
     * @type {ObjectReference}
     * @memberof DeviceResult
     */
    location?: ObjectReference;
    /**
     * 
     * @type {string}
     * @memberof DeviceResult
     */
    name?: string;
    /**
     * 
     * @type {ObjectReference}
     * @memberof DeviceResult
     */
    rack?: ObjectReference;
    /**
     * 
     * @type {ObjectReference}
     * @memberof DeviceResult
     */
    site?: ObjectReference;
    /**
     * 
     * @type {string}
     * @memberof DeviceResult
     */
    status?: string;
    /**
     * 
     * @type {Array<Tag>}
     * @memberof DeviceResult
     */
    tags?: Array<Tag>;
    /**
     * 
     * @type {Date}
     * @memberof DeviceResult
     */
    updated?: Date;
    /**
     * 
     * @type {number}
     * @memberof DeviceResult
     */
    version: number;
}

/**
 * Check if a given object implements the DeviceResult interface.
 */
export function instanceOfDeviceResult(value: object): value is DeviceResult {
    if (!('id' in value) || value['id'] === undefined) return false;
    if (!('version' in value) || value['version'] === undefined) return false;
    return true;
}

export function DeviceResultFromJSON(json: any): DeviceResult {
    return DeviceResultFromJSONTyped(json, false);
}

export function DeviceResultFromJSONTyped(json: any, ignoreDiscriminator: boolean): DeviceResult {
    if (json == null) {
        return json;
    }
    return {
        
        'created': json['created'] == null ? undefined : (new Date(json['created'])),
        'description': json['description'] == null ? undefined : json['description'],
        'deviceType': json['deviceType'] == null ? undefined : ObjectReferenceFromJSON(json['deviceType']),
        'id': json['id'],
        'location': json['location'] == null ? undefined : ObjectReferenceFromJSON(json['location']),
        'name': json['name'] == null ? undefined : json['name'],
        'rack': json['rack'] == null ? undefined : ObjectReferenceFromJSON(json['rack']),
        'site': json['site'] == null ? undefined : ObjectReferenceFromJSON(json['site']),
        'status': json['status'] == null ? undefined : json['status'],
        'tags': json['tags'] == null ? undefined : ((json['tags'] as Array<any>).map(TagFromJSON)),
        'updated': json['updated'] == null ? undefined : (new Date(json['updated'])),
        'version': json['version'],
    };
}

export function DeviceResultToJSON(json: any): DeviceResult {
    return DeviceResultToJSONTyped(json, false);
}

export function DeviceResultToJSONTyped(value?: Omit<DeviceResult, 'id'> | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'created': value['created'] == null ? undefined : ((value['created']).toISOString()),
        'description': value['description'],
        'deviceType': ObjectReferenceToJSON(value['deviceType']),
        'location': ObjectReferenceToJSON(value['location']),
        'name': value['name'],
        'rack': ObjectReferenceToJSON(value['rack']),
        'site': ObjectReferenceToJSON(value['site']),
        'status': value['status'],
        'tags': value['tags'] == null ? undefined : ((value['tags'] as Array<any>).map(TagToJSON)),
        'updated': value['updated'] == null ? undefined : ((value['updated']).toISOString()),
        'version': value['version'],
    };
}

