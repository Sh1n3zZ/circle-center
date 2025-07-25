import { post } from "../client"
import type { 
  DiffAppFiltersResponse, 
  DiffIconsResponse
} from "./types"

/**
 * Compare two appfilter.xml files and return the differences
 * @param file1 First appfilter.xml file
 * @param file2 Second appfilter.xml file
 * @returns Promise with diff results
 */
export const diffAppFilters = async (
  file1: File, 
  file2: File
): Promise<DiffAppFiltersResponse> => {
  const formData = new FormData()
  formData.append("file1", file1)
  formData.append("file2", file2)

  const response = await post<DiffAppFiltersResponse>(
    "/processor/diffappfilters",
    formData,
    { asFormData: true }
  )

  return response.data
}

/**
 * Compare local icon directory with icon_pack.xml and find missing icons
 * @param iconDir Path to icon directory
 * @param iconPack Path to icon_pack.xml file
 * @returns Promise with missing icons list
 */
export const diffIcons = async (
  iconDir: string, 
  iconPack: string
): Promise<DiffIconsResponse> => {
  const formData = new FormData()
  formData.append("icon_dir", iconDir)
  formData.append("icon_pack", iconPack)

  const response = await post<DiffIconsResponse>(
    "/processor/difficons",
    formData,
    { asFormData: true }
  )

  return response.data
}
