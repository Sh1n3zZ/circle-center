import React from "react"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card"
import type {
  Item,
  StringArray,
} from "@/api/reader/types"

// Union type for component props
interface ReaderResultItemCardProps {
  item: Item | StringArray
}

const isItem = (obj: Item | StringArray): obj is Item =>
  (obj as Item).Component !== undefined

const ReaderResultItemCard: React.FC<ReaderResultItemCardProps> = ({ item }) => {
  if (isItem(item)) {
    // Display appfilter item
    return (
      <Card className="w-full">
        <CardHeader>
          <CardTitle>{item.AppName || item.Drawable}</CardTitle>
          <CardDescription>{item.Component}</CardDescription>
        </CardHeader>
        <CardContent className="text-sm space-y-1">
          <p>
            <span className="font-medium">Package:</span> {item.PackageName}
          </p>
          <p>
            <span className="font-medium">Activity:</span> {item.ActivityName}
          </p>
          <p>
            <span className="font-medium">Drawable:</span> {item.Drawable}
          </p>
        </CardContent>
      </Card>
    )
  }

  // StringArray display (icon_pack)
  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>{item.Name}</CardTitle>
        <CardDescription>string-array</CardDescription>
      </CardHeader>
      <CardContent className="text-sm space-y-1">
        {item.Items.map((str, idx) => (
          <p key={idx}>{str}</p>
        ))}
      </CardContent>
    </Card>
  )
}

export default ReaderResultItemCard
