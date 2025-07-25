import React, { useState } from "react"
import ReaderChooseUploadFile from "@/components/reader/ReaderChooseUploadFile"
import ReaderCard from "@/components/reader/ReaderCard"
import ReaderDiffCard from "@/components/reader/ReaderDiffCard"
import ReaderMergeCard from "@/components/reader/ReaderMergeCard"
import type { ReadFileResponse } from "@/api/reader/types"
import type { DiffAppFiltersResponse } from "@/api/diff/types"

const Reader: React.FC = () => {
  const [result, setResult] = useState<ReadFileResponse | null>(null)
  const [diffResult, setDiffResult] = useState<DiffAppFiltersResponse | null>(null)
  const [mergedContent, setMergedContent] = useState<string | null>(null)
  const [file1, setFile1] = useState<File | null>(null)
  const [file2, setFile2] = useState<File | null>(null)

  const handleDiffFiles = (files: { file1: File; file2: File }) => {
    setFile1(files.file1)
    setFile2(files.file2)
  }

  return (
    <div className="p-4 space-y-6 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold">XML Reader</h1>

      <ReaderChooseUploadFile onSuccess={setResult} />

      <ReaderCard data={result} />

      <ReaderDiffCard onDiffComplete={setDiffResult} onFilesSelected={handleDiffFiles} />

      <ReaderMergeCard
        diffResult={diffResult}
        onMergeComplete={setMergedContent}
        file1={file1}
        file2={file2}
      />

      {mergedContent && (
        <div className="space-y-2">
          <h2 className="text-xl font-semibold">合并结果</h2>
          <pre className="p-4 bg-gray-50 rounded-lg overflow-x-auto">
            <code>{mergedContent}</code>
          </pre>
        </div>
      )}
    </div>
  )
}

export default Reader
