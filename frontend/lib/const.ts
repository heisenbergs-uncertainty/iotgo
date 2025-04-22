export const allowedLocations = [
    "None",
    "A",
    "B",
    "C",
    "F",
    "G"
] as const; // Define allowed values with typings

export enum LocationOptionsEnum {
    None = 0,
    A = 1,
    B = 2,
    C = 3,
    F = 4,
    G = 5,
}

export const allowedDeviceTypes = ["CNC", "Printer", "Scanner", "Other"] as const; // Define allowed values with typings