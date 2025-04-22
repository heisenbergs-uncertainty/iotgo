import {allowedDeviceTypes, allowedLocations} from "@/lib/const";

export type AllowedLocation = typeof allowedLocations[number]; // Create a type from the array

export type AllowedDeviceType = typeof allowedDeviceTypes[number]; // Create a type from the array
