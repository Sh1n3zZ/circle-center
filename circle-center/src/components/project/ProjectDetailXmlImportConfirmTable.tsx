import type { IconImportComponent } from '@/api/manager/types';
import { Input } from '@/components/ui/input';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import * as React from 'react';

type RowData = IconImportComponent & { id: number };

export default function ProjectDetailXmlImportConfirmTable({
  data,
  onChange,
}: {
  data: IconImportComponent[];
  onChange?: (next: IconImportComponent[]) => void;
}) {
  const [rows, setRows] = React.useState<RowData[]>(() =>
    data.map((d, i) => ({ ...d, id: i + 1 }))
  );

  React.useEffect(() => {
    setRows(data.map((d, i) => ({ ...d, id: i + 1 })));
  }, [data]);

  const updateRow = (idx: number, patch: Partial<IconImportComponent>) => {
    setRows(prev => {
      const next = [...prev];
      next[idx] = { ...next[idx], ...patch };
      onChange?.(next.map(({ id, ...rest }) => rest));
      return next;
    });
  };

  return (
    <div className='w-full overflow-hidden rounded-md border'>
      <Table className='w-full'>
        <TableHeader>
          <TableRow>
            <TableHead className='w-[25%] min-w-[150px]'>Name</TableHead>
            <TableHead className='w-[25%] min-w-[150px]'>Package</TableHead>
            <TableHead className='w-[25%] min-w-[150px]'>
              Component Info
            </TableHead>
            <TableHead className='w-[25%] min-w-[150px]'>Drawable</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.length ? (
            rows.map((row, idx) => (
              <TableRow key={row.id}>
                <TableCell className='p-2 w-[25%]'>
                  <Input
                    value={row.name || ''}
                    onChange={e => updateRow(idx, { name: e.target.value })}
                    className='text-xs w-full'
                    placeholder='Component name'
                  />
                </TableCell>
                <TableCell className='p-2 w-[25%]'>
                  <Input
                    value={row.pkg || ''}
                    onChange={e => updateRow(idx, { pkg: e.target.value })}
                    className='text-xs w-full'
                    placeholder='Package name'
                  />
                </TableCell>
                <TableCell className='p-2 w-[25%]'>
                  <Input
                    value={row.componentInfo || ''}
                    onChange={e =>
                      updateRow(idx, { componentInfo: e.target.value })
                    }
                    className='text-xs w-full'
                    placeholder='Component info'
                  />
                </TableCell>
                <TableCell className='p-2 w-[25%]'>
                  <Input
                    value={row.drawable || ''}
                    onChange={e => updateRow(idx, { drawable: e.target.value })}
                    className='text-xs w-full'
                    placeholder='Drawable resource'
                  />
                </TableCell>
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell
                colSpan={4}
                className='text-center text-muted-foreground'
              >
                No data available
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
