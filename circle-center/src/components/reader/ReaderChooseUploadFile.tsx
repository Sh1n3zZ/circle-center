import React, { useRef, useState } from "react"
import { Button } from "@/components/ui/button"
import { uploadReaderFile } from "@/api/reader/readfile"
import type { ReadFileResponse, ReaderFileType } from "@/api/reader/types"
import { toast } from "sonner"

interface ReaderChooseUploadFileProps {
  /** Callback when backend returns parsed data successfully */
  onSuccess: (data: ReadFileResponse) => void
  /** Optional: specify which type to parse. Defaults to "appfilter" */
  type?: ReaderFileType
}

const ReaderChooseUploadFile: React.FC<ReaderChooseUploadFileProps> = ({
  onSuccess,
  type = "appfilter",
}) => {
  const fileInputRef = useRef<HTMLInputElement | null>(null)
  const [loading, setLoading] = useState(false)

  const handleButtonClick = () => {
    fileInputRef.current?.click()
  }

  const handleFileChange: React.ChangeEventHandler<HTMLInputElement> = async (
    event,
  ) => {
    const file = event.target.files?.[0]
    if (!file) return

    try {
      setLoading(true)
      const data = await uploadReaderFile(file, type)
      toast.success("File parsed successfully")
      onSuccess(data)
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to parse the xml file",
      )
    } finally {
      setLoading(false)
      // Reset input value so that selecting the same file again triggers change
      event.target.value = ""
    }
  }

  return (
    <div className="flex items-center gap-4">
      <input
        ref={fileInputRef}
        type="file"
        accept=".xml"
        className="hidden"
        onChange={handleFileChange}
      />
      <Button disabled={loading} onClick={handleButtonClick}>
        {loading ? "Uploading..." : "Choose XML File"}
      </Button>
    </div>
  )
}

export default ReaderChooseUploadFile
