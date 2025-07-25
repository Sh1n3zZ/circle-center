import React from "react"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"

interface DiffItemProps {
  component: string
  drawable: string
  PackageName: string
  ActivityName: string
  AppName?: string
  type: "only-first" | "only-second" | "common"
}

const ReaderDiffItem: React.FC<DiffItemProps> = ({
  component,
  drawable,
  PackageName,
  ActivityName,
  AppName,
  type,
}) => {
  const getTypeColor = () => {
    switch (type) {
      case "only-first":
        return "bg-green-100 border-green-300 text-green-800"
      case "only-second":
        return "bg-red-100 border-red-300 text-red-800"
      case "common":
        return "bg-gray-100 border-gray-300 text-gray-800"
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
      case "common":
        return "共同项"
      default:
        return "未知"
    }
  }

  return (
    <Card className={`border-2 ${getTypeColor()}`}>
      <CardContent className="p-4">
        <div className="flex items-center justify-between mb-2">
          <Badge variant="outline" className={getTypeColor()}>
            {getTypeLabel()}
          </Badge>
        </div>
        
        <div className="space-y-2">
          {AppName && (
            <div className="flex items-center space-x-2">
              <span className="font-medium text-sm">应用名称:</span>
              <span className="text-sm">{AppName}</span>
            </div>
          )}
          
          <div className="flex items-center space-x-2">
            <span className="font-medium text-sm">包名:</span>
            <span className="text-sm font-mono">{PackageName}</span>
          </div>
          
          <div className="flex items-center space-x-2">
            <span className="font-medium text-sm">活动名:</span>
            <span className="text-sm font-mono">{ActivityName}</span>
          </div>
          
          <div className="flex items-center space-x-2">
            <span className="font-medium text-sm">组件:</span>
            <span className="text-sm font-mono text-xs break-all">{component}</span>
          </div>
          
          <div className="flex items-center space-x-2">
            <span className="font-medium text-sm">图标:</span>
            <span className="text-sm font-mono">{drawable}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

export default ReaderDiffItem
