import client from "../client"
import type {
  ReaderFileType,
  ReadFileResponse,
} from "./types"

/**
 * Uploads an xml file to backend /reader/readfile endpoint.
 * The backend accepts multipart/form-data with fields:
 * - file: the xml file
 * - type: optional, either "appfilter" or "icon_pack"
 */
export async function uploadReaderFile(
  file: File,
  type: ReaderFileType = "appfilter",
): Promise<ReadFileResponse> {
  const formData = new FormData()
  formData.append("file", file)
  if (type) {
    formData.append("type", type)
  }

  const { data } = await client.post<ReadFileResponse>(
    "/reader/readfile",
    formData,
    {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    },
  )

  return data
}
