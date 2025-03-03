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
 * @interface LocationResult
 */
export interface LocationResult {
    /**
     * 
     * @type {Date}
     * @memberof LocationResult
     */
    readonly created?: Date;
    /**
     * 
     * @type {string}
     * @memberof LocationResult
     */
    description?: string;
    /**
     * 
     * @type {string}
     * @memberof LocationResult
     */
    readonly id: string;
    /**
     * 
     * @type {string}
     * @memberof LocationResult
     */
    name: string;
    /**
     * 
     * @type {ObjectReference}
     * @memberof LocationResult
     */
    parent?: ObjectReference;
    /**
     * 
     * @type {ObjectReference}
     * @memberof LocationResult
     */
    site?: ObjectReference;
    /**
     * 
     * @type {Array<Tag>}
     * @memberof LocationResult
     */
    tags?: Array<Tag>;
    /**
     * 
     * @type {Date}
     * @memberof LocationResult
     */
    readonly updated?: Date;
    /**
     * 
     * @type {number}
     * @memberof LocationResult
     */
    version: number;
}

/**
 * Check if a given object implements the LocationResult interface.
 */
export function instanceOfLocationResult(value: object): value is LocationResult {
    if (!('id' in value) || value['id'] === undefined) return false;
    if (!('name' in value) || value['name'] === undefined) return false;
    if (!('version' in value) || value['version'] === undefined) return false;
    return true;
}

export function LocationResultFromJSON(json: any): LocationResult {
    return LocationResultFromJSONTyped(json, false);
}

export function LocationResultFromJSONTyped(json: any, ignoreDiscriminator: boolean): LocationResult {
    if (json == null) {
        return json;
    }
    return {
        
        'created': json['created'] == null ? undefined : (new Date(json['created'])),
        'description': json['description'] == null ? undefined : json['description'],
        'id': json['id'],
        'name': json['name'],
        'parent': json['parent'] == null ? undefined : ObjectReferenceFromJSON(json['parent']),
        'site': json['site'] == null ? undefined : ObjectReferenceFromJSON(json['site']),
        'tags': json['tags'] == null ? undefined : ((json['tags'] as Array<any>).map(TagFromJSON)),
        'updated': json['updated'] == null ? undefined : (new Date(json['updated'])),
        'version': json['version'],
    };
}

export function LocationResultToJSON(json: any): LocationResult {
    return LocationResultToJSONTyped(json, false);
}

export function LocationResultToJSONTyped(value?: Omit<LocationResult, 'created'|'id'|'updated'> | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'description': value['description'],
        'name': value['name'],
        'parent': ObjectReferenceToJSON(value['parent']),
        'site': ObjectReferenceToJSON(value['site']),
        'tags': value['tags'] == null ? undefined : ((value['tags'] as Array<any>).map(TagToJSON)),
        'version': value['version'],
    };
}

