import React from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import type { ReadFileResponse, Item, StringArray } from "@/api/reader/types"
import ReaderResultItemCard from "./ReaderResultItemCard"

interface ReaderCardProps {
  data: ReadFileResponse | null
}

const ReaderCard: React.FC<ReaderCardProps> = ({ data }) => {
  if (!data) {
    return (
      <Card className="w-full">
        <CardHeader>
          <CardTitle>No data yet. Upload a file to see results.</CardTitle>
        </CardHeader>
      </Card>
    )
  }

  // Determine whether items or icon_pack
  const items: (Item | StringArray)[] =
    "items" in data ? data.items : data.icon_pack.Arrays

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Parsed Result</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {items.length === 0 ? (
          <p>No items found.</p>
        ) : (
          items.map((it, idx) => (
            <ReaderResultItemCard key={idx} item={it} />
          ))
        )}
      </CardContent>
    </Card>
  )
}

export default ReaderCard
