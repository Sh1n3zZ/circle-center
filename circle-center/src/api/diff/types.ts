export interface DiffItem {
  component: string
  drawable: string
  PackageName: string
  ActivityName: string
  AppName?: string
}

export interface DiffAppFiltersResponse {
  only_in_first: DiffItem[]
  only_in_second: DiffItem[]
  common: DiffItem[]
  summary: {
    first_count: number
    second_count: number
    only_first_count: number
    only_second_count: number
    common_count: number
  }
}

export interface DiffIconsResponse {
  missing_icons: string[]
  count: number
}

export interface DiffAppFiltersRequest {
  file1: File
  file2: File
}

export interface DiffIconsRequest {
  icon_dir: string
  icon_pack: string
}
