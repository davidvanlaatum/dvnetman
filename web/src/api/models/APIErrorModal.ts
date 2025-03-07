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
import type { ErrorMessage } from './ErrorMessage';
import {
    ErrorMessageFromJSON,
    ErrorMessageFromJSONTyped,
    ErrorMessageToJSON,
    ErrorMessageToJSONTyped,
} from './ErrorMessage';

/**
 * 
 * @export
 * @interface APIErrorModal
 */
export interface APIErrorModal {
    /**
     * 
     * @type {Array<ErrorMessage>}
     * @memberof APIErrorModal
     */
    errors?: Array<ErrorMessage>;
}

/**
 * Check if a given object implements the APIErrorModal interface.
 */
export function instanceOfAPIErrorModal(value: object): value is APIErrorModal {
    return true;
}

export function APIErrorModalFromJSON(json: any): APIErrorModal {
    return APIErrorModalFromJSONTyped(json, false);
}

export function APIErrorModalFromJSONTyped(json: any, ignoreDiscriminator: boolean): APIErrorModal {
    if (json == null) {
        return json;
    }
    return {
        
        'errors': json['errors'] == null ? undefined : ((json['errors'] as Array<any>).map(ErrorMessageFromJSON)),
    };
}

export function APIErrorModalToJSON(json: any): APIErrorModal {
    return APIErrorModalToJSONTyped(json, false);
}

export function APIErrorModalToJSONTyped(value?: APIErrorModal | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'errors': value['errors'] == null ? undefined : ((value['errors'] as Array<any>).map(ErrorMessageToJSON)),
    };
}

