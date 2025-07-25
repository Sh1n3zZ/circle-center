import React, { useState } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Badge } from "@/components/ui/badge"
import ReaderDiffItem from "./ReaderDiffItem"
import { diffAppFilters } from "@/api/diff/diff"
import type { DiffAppFiltersResponse } from "@/api/diff/types"

interface ReaderDiffCardProps {
  onDiffComplete?: (result: DiffAppFiltersResponse) => void
}

const ReaderDiffCard: React.FC<ReaderDiffCardProps> = ({ onDiffComplete }) => {
  const [file1, setFile1] = useState<File | null>(null)
  const [file2, setFile2] = useState<File | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [diffResult, setDiffResult] = useState<DiffAppFiltersResponse | null>(null)
  const [error, setError] = useState<string | null>(null)

  const handleFile1Change = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    setFile1(file || null)
  }

  const handleFile2Change = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    setFile2(file || null)
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
          <div className="space-y-4">
            {/* Summary */}
            <div className="grid grid-cols-2 md:grid-cols-5 gap-4 p-4 bg-gray-50 rounded-lg">
              <div className="text-center">
                <div className="text-2xl font-bold text-green-600">
                  {diffResult.summary.only_first_count}
                </div>
                <div className="text-sm text-gray-600">仅在第一个文件中</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-red-600">
                  {diffResult.summary.only_second_count}
                </div>
                <div className="text-sm text-gray-600">仅在第二个文件中</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-gray-600">
                  {diffResult.summary.common_count}
                </div>
                <div className="text-sm text-gray-600">共同项</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-blue-600">
                  {diffResult.summary.first_count}
                </div>
                <div className="text-sm text-gray-600">第一个文件总数</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-purple-600">
                  {diffResult.summary.second_count}
                </div>
                <div className="text-sm text-gray-600">第二个文件总数</div>
              </div>
            </div>

            {/* Detailed Results */}
            <Tabs defaultValue="only-first" className="w-full">
              <TabsList className="grid w-full grid-cols-3">
                <TabsTrigger value="only-first">
                  仅在第一个文件中 ({diffResult?.only_in_first?.length || 0})
                </TabsTrigger>
                <TabsTrigger value="only-second">
                  仅在第二个文件中 ({diffResult?.only_in_second?.length || 0})
                </TabsTrigger>
                <TabsTrigger value="common">
                  共同项 ({diffResult?.common?.length || 0})
                </TabsTrigger>
              </TabsList>
              
              <TabsContent value="only-first" className="space-y-4">
                {!diffResult?.only_in_first || diffResult.only_in_first.length === 0 ? (
                  <p className="text-center text-gray-500 py-8">没有仅在第一个文件中的项</p>
                ) : (
                  <div className="grid gap-4">
                    {diffResult.only_in_first.map((item, index) => (
                      <ReaderDiffItem
                        key={`only-first-${index}`}
                        {...item}
                        type="only-first"
                      />
                    ))}
                  </div>
                )}
              </TabsContent>
              
              <TabsContent value="only-second" className="space-y-4">
                {!diffResult?.only_in_second || diffResult.only_in_second.length === 0 ? (
                  <p className="text-center text-gray-500 py-8">没有仅在第二个文件中的项</p>
                ) : (
                  <div className="grid gap-4">
                    {diffResult.only_in_second.map((item, index) => (
                      <ReaderDiffItem
                        key={`only-second-${index}`}
                        {...item}
                        type="only-second"
                      />
                    ))}
                  </div>
                )}
              </TabsContent>
              
              <TabsContent value="common" className="space-y-4">
                {!diffResult?.common || diffResult.common.length === 0 ? (
                  <p className="text-center text-gray-500 py-8">没有共同项</p>
                ) : (
                  <div className="grid gap-4">
                    {diffResult.common.map((item, index) => (
                      <ReaderDiffItem
                        key={`common-${index}`}
                        {...item}
                        type="common"
                      />
                    ))}
                  </div>
                )}
              </TabsContent>
            </Tabs>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

export default ReaderDiffCard
