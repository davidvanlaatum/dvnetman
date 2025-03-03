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
import type { ObjectReference } from './ObjectReference';
import {
    ObjectReferenceFromJSON,
    ObjectReferenceFromJSONTyped,
    ObjectReferenceToJSON,
    ObjectReferenceToJSONTyped,
} from './ObjectReference';
import type { DevicePort } from './DevicePort';
import {
    DevicePortFromJSON,
    DevicePortFromJSONTyped,
    DevicePortToJSON,
    DevicePortToJSONTyped,
} from './DevicePort';
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
 * @interface Device
 */
export interface Device {
    /**
     * 
     * @type {string}
     * @memberof Device
     */
    assetTag?: string;
    /**
     * 
     * @type {Date}
     * @memberof Device
     */
    readonly created?: Date;
    /**
     * 
     * @type {string}
     * @memberof Device
     */
    description?: string;
    /**
     * 
     * @type {ObjectReference}
     * @memberof Device
     */
    deviceType?: ObjectReference;
    /**
     * 
     * @type {string}
     * @memberof Device
     */
    readonly id: string;
    /**
     * 
     * @type {ObjectReference}
     * @memberof Device
     */
    location?: ObjectReference;
    /**
     * 
     * @type {string}
     * @memberof Device
     */
    name?: string;
    /**
     * 
     * @type {Array<DevicePort>}
     * @memberof Device
     */
    ports?: Array<DevicePort>;
    /**
     * 
     * @type {ObjectReference}
     * @memberof Device
     */
    rack?: ObjectReference;
    /**
     * 
     * @type {string}
     * @memberof Device
     */
    rackFace?: DeviceRackFaceEnum;
    /**
     * 
     * @type {number}
     * @memberof Device
     */
    rackPosition?: number;
    /**
     * 
     * @type {string}
     * @memberof Device
     */
    serial?: string;
    /**
     * 
     * @type {ObjectReference}
     * @memberof Device
     */
    site?: ObjectReference;
    /**
     * 
     * @type {string}
     * @memberof Device
     */
    status?: string;
    /**
     * 
     * @type {Array<Tag>}
     * @memberof Device
     */
    tags?: Array<Tag>;
    /**
     * 
     * @type {Date}
     * @memberof Device
     */
    readonly updated?: Date;
    /**
     * 
     * @type {number}
     * @memberof Device
     */
    version: number;
}


/**
 * @export
 */
export const DeviceRackFaceEnum = {
    Front: 'front',
    Rear: 'rear'
} as const;
export type DeviceRackFaceEnum = typeof DeviceRackFaceEnum[keyof typeof DeviceRackFaceEnum];


/**
 * Check if a given object implements the Device interface.
 */
export function instanceOfDevice(value: object): value is Device {
    if (!('id' in value) || value['id'] === undefined) return false;
    if (!('version' in value) || value['version'] === undefined) return false;
    return true;
}

export function DeviceFromJSON(json: any): Device {
    return DeviceFromJSONTyped(json, false);
}

export function DeviceFromJSONTyped(json: any, ignoreDiscriminator: boolean): Device {
    if (json == null) {
        return json;
    }
    return {
        
        'assetTag': json['assetTag'] == null ? undefined : json['assetTag'],
        'created': json['created'] == null ? undefined : (new Date(json['created'])),
        'description': json['description'] == null ? undefined : json['description'],
        'deviceType': json['deviceType'] == null ? undefined : ObjectReferenceFromJSON(json['deviceType']),
        'id': json['id'],
        'location': json['location'] == null ? undefined : ObjectReferenceFromJSON(json['location']),
        'name': json['name'] == null ? undefined : json['name'],
        'ports': json['ports'] == null ? undefined : ((json['ports'] as Array<any>).map(DevicePortFromJSON)),
        'rack': json['rack'] == null ? undefined : ObjectReferenceFromJSON(json['rack']),
        'rackFace': json['rackFace'] == null ? undefined : json['rackFace'],
        'rackPosition': json['rackPosition'] == null ? undefined : json['rackPosition'],
        'serial': json['serial'] == null ? undefined : json['serial'],
        'site': json['site'] == null ? undefined : ObjectReferenceFromJSON(json['site']),
        'status': json['status'] == null ? undefined : json['status'],
        'tags': json['tags'] == null ? undefined : ((json['tags'] as Array<any>).map(TagFromJSON)),
        'updated': json['updated'] == null ? undefined : (new Date(json['updated'])),
        'version': json['version'],
    };
}

export function DeviceToJSON(json: any): Device {
    return DeviceToJSONTyped(json, false);
}

export function DeviceToJSONTyped(value?: Omit<Device, 'created'|'id'|'updated'> | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'assetTag': value['assetTag'],
        'description': value['description'],
        'deviceType': ObjectReferenceToJSON(value['deviceType']),
        'location': ObjectReferenceToJSON(value['location']),
        'name': value['name'],
        'ports': value['ports'] == null ? undefined : ((value['ports'] as Array<any>).map(DevicePortToJSON)),
        'rack': ObjectReferenceToJSON(value['rack']),
        'rackFace': value['rackFace'],
        'rackPosition': value['rackPosition'],
        'serial': value['serial'],
        'site': ObjectReferenceToJSON(value['site']),
        'status': value['status'],
        'tags': value['tags'] == null ? undefined : ((value['tags'] as Array<any>).map(TagToJSON)),
        'version': value['version'],
    };
}

