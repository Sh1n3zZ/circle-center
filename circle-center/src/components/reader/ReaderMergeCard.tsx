import React, { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import { Badge } from "@/components/ui/badge"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import ReaderMergeItem from "./ReaderMergeItem"
import type { DiffAppFiltersResponse } from "@/api/diff/types"
import { mergeAppFilters, downloadMergedXml } from "@/api/merge/merge"
import type { MergeResponse } from "@/api/merge/types"

interface ReaderMergeCardProps {
  diffResult?: DiffAppFiltersResponse | null
  onMergeComplete?: (content: string) => void
  file1?: File | null
  file2?: File | null
}

const ReaderMergeCard: React.FC<ReaderMergeCardProps> = ({
  diffResult,
  onMergeComplete,
  file1,
  file2,
}) => {
  const [selectedItems, setSelectedItems] = useState<Set<string>>(new Set())
  const [mergeIntoFirst, setMergeIntoFirst] = useState(true)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Reset selections when diff result changes
  useEffect(() => {
    setSelectedItems(new Set())
  }, [diffResult])

  const handleSelectAll = (type: "only-first" | "only-second") => {
    const items = type === "only-first" ? diffResult?.only_in_first : diffResult?.only_in_second
    if (!items) return

    const newSelected = new Set(selectedItems)
    items.forEach(item => newSelected.add(item.component))
    setSelectedItems(newSelected)
  }

  const handleDeselectAll = (type: "only-first" | "only-second") => {
    const items = type === "only-first" ? diffResult?.only_in_first : diffResult?.only_in_second
    if (!items) return

    const newSelected = new Set(selectedItems)
    items.forEach(item => newSelected.delete(item.component))
    setSelectedItems(newSelected)
  }

  const handleItemSelect = (component: string, selected: boolean) => {
    const newSelected = new Set(selectedItems)
    if (selected) {
      newSelected.add(component)
    } else {
      newSelected.delete(component)
    }
    setSelectedItems(newSelected)
  }

  const handleMerge = async () => {
    if (!diffResult) return

    setIsLoading(true)
    setError(null)

    try {
      const result = await mergeAppFilters(
        file1!,
        file2!,
        {
          components: Array.from(selectedItems),
          merge_into_first: mergeIntoFirst,
        }
      )

      onMergeComplete?.(result.content)

      // Offer download of merged file
      downloadMergedXml(result.content, "merged_appfilter.xml")
    } catch (err) {
      setError(err instanceof Error ? err.message : "合并失败")
    } finally {
      setIsLoading(false)
    }
  }

  if (!diffResult) return null

  const { only_in_first, only_in_second } = diffResult

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span>合并差异</span>
          <div className="flex items-center space-x-2">
            <span className="text-sm">合并方向:</span>
            <div className="flex items-center space-x-2">
              <Switch
                checked={mergeIntoFirst}
                onCheckedChange={setMergeIntoFirst}
                id="merge-direction"
              />
              <Label htmlFor="merge-direction" className="text-sm">
                {mergeIntoFirst ? "合并到第一个文件" : "合并到第二个文件"}
              </Label>
            </div>
          </div>
        </CardTitle>
      </CardHeader>

      <CardContent className="space-y-6">
        {error && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        <Tabs defaultValue="only-first" className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="only-first">
              仅在第一个文件中 ({only_in_first.length})
            </TabsTrigger>
            <TabsTrigger value="only-second">
              仅在第二个文件中 ({only_in_second.length})
            </TabsTrigger>
          </TabsList>

          <TabsContent value="only-first" className="space-y-4">
            <div className="flex justify-between items-center">
              <Badge variant="outline">
                已选择 {only_in_first.filter(item => selectedItems.has(item.component)).length} 项
              </Badge>
              <div className="space-x-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleSelectAll("only-first")}
                >
                  全选
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleDeselectAll("only-first")}
                >
                  取消全选
                </Button>
              </div>
            </div>

            <div className="grid gap-4">
              {only_in_first.map((item) => (
                <ReaderMergeItem
                  key={item.component}
                  item={item}
                  type="only-first"
                  selected={selectedItems.has(item.component)}
                  onSelect={(selected) => handleItemSelect(item.component, selected)}
                />
              ))}
            </div>
          </TabsContent>

          <TabsContent value="only-second" className="space-y-4">
            <div className="flex justify-between items-center">
              <Badge variant="outline">
                已选择 {only_in_second.filter(item => selectedItems.has(item.component)).length} 项
              </Badge>
              <div className="space-x-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleSelectAll("only-second")}
                >
                  全选
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleDeselectAll("only-second")}
                >
                  取消全选
                </Button>
              </div>
            </div>

            <div className="grid gap-4">
              {only_in_second.map((item) => (
                <ReaderMergeItem
                  key={item.component}
                  item={item}
                  type="only-second"
                  selected={selectedItems.has(item.component)}
                  onSelect={(selected) => handleItemSelect(item.component, selected)}
                />
              ))}
            </div>
          </TabsContent>
        </Tabs>

        <Button
          className="w-full"
          disabled={selectedItems.size === 0 || isLoading}
          onClick={handleMerge}
        >
          {isLoading ? "合并中..." : `合并选中的 ${selectedItems.size} 项`}
        </Button>
      </CardContent>
    </Card>
  )
}

export default ReaderMergeCard
