import { post } from "../client"
import type { MergeRequest, MergeResponse } from "./types"

/**
 * Merge two appfilter.xml files according to the specified options.
 * @param file1 First appfilter.xml file
 * @param file2 Second appfilter.xml file
 * @param options Merge options (selected components and merge direction)
 * @returns Promise with merge result and merged content
 */
export const mergeAppFilters = async (
  file1: File,
  file2: File,
  options: MergeRequest
): Promise<MergeResponse> => {
  const formData = new FormData()
  formData.append("file1", file1)
  formData.append("file2", file2)
  formData.append("merge_into_first", String(options.merge_into_first))
  formData.append("components", JSON.stringify(options.components))

  const response = await post<MergeResponse>(
    "/processor/mergeappfilters",
    formData,
    { asFormData: true }
  )

  return response.data
}

/**
 * Download the merged XML content as a file.
 * @param content The merged XML content
 * @param filename The name for the downloaded file
 */
export const downloadMergedXml = (content: string, filename = "merged.xml") => {
  const blob = new Blob([content], { type: "application/xml" })
  const url = URL.createObjectURL(blob)
  const link = document.createElement("a")
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}
