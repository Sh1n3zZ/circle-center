import type { DiffItem } from "../diff/types"

export interface MergeRequest {
  // Selected components to merge (empty means merge all)
  components: string[]
  // If true, merge into first file; if false, merge into second file
  merge_into_first: boolean
}

export interface MergeResult {
  // Number of items merged from each file
  items_merged: number
  // Any items that failed to merge (e.g., duplicates)
  failed_items: DiffItem[]
  // Final item count in the output file
  total_items: number
}

export interface MergeResponse {
  // The merge operation result
  result: MergeResult
  // The merged XML content
  content: string
}
