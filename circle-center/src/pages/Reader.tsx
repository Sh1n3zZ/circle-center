import React, { useState } from "react"
import ReaderChooseUploadFile from "@/components/reader/ReaderChooseUploadFile"
import ReaderCard from "@/components/reader/ReaderCard"
import type { ReadFileResponse } from "@/api/reader/types"

const Reader: React.FC = () => {
  const [result, setResult] = useState<ReadFileResponse | null>(null)

  return (
    <div className="p-4 space-y-6 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold">XML Reader</h1>

      <ReaderChooseUploadFile onSuccess={setResult} />

      <ReaderCard data={result} />
    </div>
  )
}

export default Reader
