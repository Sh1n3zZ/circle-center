import React from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Card, CardContent } from "@/components/ui/card"
import ReaderDiffItem from "./ReaderDiffItem"
import type { DiffAppFiltersResponse } from "@/api/diff/types"

interface ReaderDiffInfoProps {
  diffResult: DiffAppFiltersResponse
}

const ReaderDiffInfo: React.FC<ReaderDiffInfoProps> = ({ diffResult }) => {
  return (
    <div className="w-full flex flex-col gap-6">
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
        <TabsList className="w-full flex">
          <TabsTrigger className="flex-1" value="only-first">
            仅在第一个文件中 ({diffResult?.only_in_first?.length || 0})
          </TabsTrigger>
          <TabsTrigger className="flex-1" value="only-second">
            仅在第二个文件中 ({diffResult?.only_in_second?.length || 0})
          </TabsTrigger>
          <TabsTrigger className="flex-1" value="common">
            共同项 ({diffResult?.common?.length || 0})
          </TabsTrigger>
        </TabsList>
        
        <TabsContent value="only-first">
          <Card>
            <CardContent>
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
            </CardContent>
          </Card>
        </TabsContent>
        
        <TabsContent value="only-second">
          <Card>
            <CardContent>
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
            </CardContent>
          </Card>
        </TabsContent>
        
        <TabsContent value="common">
          <Card>
            <CardContent>
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
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

export default ReaderDiffInfo
