export type ReaderFileType = "appfilter" | "icon_pack"

// Represents one <item> element parsed from appfilter.xml coming from backend.
export interface Item {
  Component: string
  Drawable: string
  AppName: string
  PackageName: string
  ActivityName: string
}

// The successful response shape when uploading appfilter.xml
export interface ItemsResponse {
  items: Item[]
}

// Represents <string-array> element inside icon_pack.xml
export interface StringArray {
  Name: string
  Items: string[]
}

// Root structure returned by backend for icon_pack.xml
export interface IconPackResources {
  Arrays: StringArray[]
}

// The successful response shape when uploading icon_pack.xml
export interface IconPackResponse {
  icon_pack: IconPackResources
}

// Union type of all possible backend responses
export type ReadFileResponse = ItemsResponse | IconPackResponse
