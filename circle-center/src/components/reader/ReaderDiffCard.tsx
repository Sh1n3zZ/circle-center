import React, { useState } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"
import { diffAppFilters } from "@/api/diff/diff"
import type { DiffAppFiltersResponse } from "@/api/diff/types"
import ReaderDiffInfo from "./ReaderDiffInfo"

interface ReaderDiffCardProps {
  onDiffComplete?: (result: DiffAppFiltersResponse) => void
  onFilesSelected?: (files: { file1: File; file2: File }) => void
}

const ReaderDiffCard: React.FC<ReaderDiffCardProps> = ({
  onDiffComplete,
  onFilesSelected,
}) => {
  const [file1, setFile1] = useState<File | null>(null)
  const [file2, setFile2] = useState<File | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [diffResult, setDiffResult] = useState<DiffAppFiltersResponse | null>(null)
  const [error, setError] = useState<string | null>(null)

  const handleFile1Change = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    setFile1(file || null)
    if (file && file2) {
      onFilesSelected?.({ file1: file, file2 })
    }
  }

  const handleFile2Change = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    setFile2(file || null)
    if (file1 && file) {
      onFilesSelected?.({ file1, file2: file })
    }
  }

  const handleDiff = async () => {
    if (!file1 || !file2) {
      setError("请选择两个文件进行比较")
      return
    }

    setIsLoading(true)
    setError(null)

    try {
      const result = await diffAppFilters(file1, file2)
      setDiffResult(result)
      onDiffComplete?.(result)
    } catch (err) {
      setError(err instanceof Error ? err.message : "比较失败")
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="flex items-center space-x-2">
          <span>文件比较</span>
          {diffResult && (
            <Badge variant="secondary">
              总计: {diffResult.summary.first_count + diffResult.summary.second_count}
            </Badge>
          )}
        </CardTitle>
      </CardHeader>
      
      <CardContent className="space-y-6">
        {/* File Upload Section */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label htmlFor="file1">第一个文件 (主文件)</Label>
            <Input
              id="file1"
              type="file"
              accept=".xml"
              onChange={handleFile1Change}
              className="cursor-pointer"
            />
            {file1 && (
              <p className="text-sm text-gray-600">{file1.name}</p>
            )}
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="file2">第二个文件</Label>
            <Input
              id="file2"
              type="file"
              accept=".xml"
              onChange={handleFile2Change}
              className="cursor-pointer"
            />
            {file2 && (
              <p className="text-sm text-gray-600">{file2.name}</p>
            )}
          </div>
        </div>

        {/* Diff Button */}
        <Button 
          onClick={handleDiff} 
          disabled={!file1 || !file2 || isLoading}
          className="w-full"
        >
          {isLoading ? "比较中..." : "开始比较"}
        </Button>

        {/* Error Display */}
        {error && (
          <div className="p-4 bg-red-50 border border-red-200 rounded-md">
            <p className="text-red-800 text-sm">{error}</p>
          </div>
        )}

        {/* Diff Results */}
        {diffResult && (
          <ReaderDiffInfo diffResult={diffResult} />
        )}
      </CardContent>
    </Card>
  )
}

export default ReaderDiffCard
