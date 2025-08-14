import * as React from "react";
import { Input } from "@/components/ui/input";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import type { IconImportComponent } from "@/api/manager/types";

type RowData = IconImportComponent & { id: number };

export default function ProjectDetailXmlImportConfirmTable({ data, onChange }: { data: IconImportComponent[]; onChange?: (next: IconImportComponent[]) => void }) {
  const [rows, setRows] = React.useState<RowData[]>(() => data.map((d, i) => ({ ...d, id: i + 1 })));

  React.useEffect(() => {
    setRows(data.map((d, i) => ({ ...d, id: i + 1 })));
  }, [data]);

  const updateRow = (idx: number, patch: Partial<IconImportComponent>) => {
    setRows((prev) => {
      const next = [...prev];
      next[idx] = { ...next[idx], ...patch };
      onChange?.(next.map(({ id, ...rest }) => rest));
      return next;
    });
  };

  return (
    <div className="overflow-hidden rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Package</TableHead>
            <TableHead>ComponentInfo</TableHead>
            <TableHead>Drawable</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.length ? (
            rows.map((row, idx) => (
              <TableRow key={row.id}>
                <TableCell>
                  <Input value={row.name || ""} onChange={(e) => updateRow(idx, { name: e.target.value })} />
                </TableCell>
                <TableCell>
                  <Input value={row.pkg || ""} onChange={(e) => updateRow(idx, { pkg: e.target.value })} />
                </TableCell>
                <TableCell>
                  <Input value={row.componentInfo || ""} onChange={(e) => updateRow(idx, { componentInfo: e.target.value })} />
                </TableCell>
                <TableCell>
                  <Input value={row.drawable || ""} onChange={(e) => updateRow(idx, { drawable: e.target.value })} />
                </TableCell>
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={4} className="text-center">No data</TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}


