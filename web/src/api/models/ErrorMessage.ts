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
 * @interface ErrorMessage
 */
export interface ErrorMessage {
    /**
     * 
     * @type {string}
     * @memberof ErrorMessage
     */
    code: string;
    /**
     * 
     * @type {string}
     * @memberof ErrorMessage
     */
    message: string;
}

/**
 * Check if a given object implements the ErrorMessage interface.
 */
export function instanceOfErrorMessage(value: object): value is ErrorMessage {
    if (!('code' in value) || value['code'] === undefined) return false;
    if (!('message' in value) || value['message'] === undefined) return false;
    return true;
}

export function ErrorMessageFromJSON(json: any): ErrorMessage {
    return ErrorMessageFromJSONTyped(json, false);
}

export function ErrorMessageFromJSONTyped(json: any, ignoreDiscriminator: boolean): ErrorMessage {
    if (json == null) {
        return json;
    }
    return {
        
        'code': json['code'],
        'message': json['message'],
    };
}

export function ErrorMessageToJSON(json: any): ErrorMessage {
    return ErrorMessageToJSONTyped(json, false);
}

export function ErrorMessageToJSONTyped(value?: ErrorMessage | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'code': value['code'],
        'message': value['message'],
    };
}

