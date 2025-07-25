import React from "react"
import { Card, CardContent } from "@/components/ui/card"
import { Checkbox } from "@/components/ui/checkbox"
import { Badge } from "@/components/ui/badge"
import type { DiffItem } from "@/api/diff/types"

interface ReaderMergeItemProps {
  item: DiffItem
  selected: boolean
  onSelect: (selected: boolean) => void
  type: "only-first" | "only-second"
}

const ReaderMergeItem: React.FC<ReaderMergeItemProps> = ({
  item,
  selected,
  onSelect,
  type,
}) => {
  const getTypeColor = () => {
    switch (type) {
      case "only-first":
        return "bg-green-100 border-green-300 text-green-800"
      case "only-second":
        return "bg-red-100 border-red-300 text-red-800"
      default:
        return "bg-gray-100 border-gray-300 text-gray-800"
    }
  }

  const getTypeLabel = () => {
    switch (type) {
      case "only-first":
        return "仅在第一个文件中"
      case "only-second":
        return "仅在第二个文件中"
      default:
        return "未知"
    }
  }

  return (
    <Card className={`border-2 ${getTypeColor()}`}>
      <CardContent className="p-4">
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center space-x-2">
            <Checkbox
              checked={selected}
              onCheckedChange={onSelect}
              id={`merge-item-${item.component}`}
            />
            <Badge variant="outline" className={getTypeColor()}>
              {getTypeLabel()}
            </Badge>
          </div>
        </div>
        
        <div className="space-y-2">
          {item.AppName && (
            <div className="flex items-center space-x-2">
              <span className="font-medium text-sm">应用名称:</span>
              <span className="text-sm">{item.AppName}</span>
            </div>
          )}
          
          <div className="flex items-center space-x-2">
            <span className="font-medium text-sm">包名:</span>
            <span className="text-sm font-mono">{item.PackageName}</span>
          </div>
          
          <div className="flex items-center space-x-2">
            <span className="font-medium text-sm">活动名:</span>
            <span className="text-sm font-mono">{item.ActivityName}</span>
          </div>
          
          <div className="flex items-center space-x-2">
            <span className="font-medium text-sm">组件:</span>
            <span className="text-sm font-mono text-xs break-all">{item.component}</span>
          </div>
          
          <div className="flex items-center space-x-2">
            <span className="font-medium text-sm">图标:</span>
            <span className="text-sm font-mono">{item.drawable}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

export default ReaderMergeItem
