import type { IconModel } from '@/api/manager/types';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { LoadingSpinner } from '@/components/ui/loading-spinner';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  useReactTable,
  type ColumnDef,
  type ColumnFiltersState,
  type SortingState,
} from '@tanstack/react-table';
import { ChevronDown } from 'lucide-react';
import * as React from 'react';
import ProjectDetailIconDrawableView from './ProjectDetailIconDrawableView';
import ProjectDetailIconTableActions from './ProjectDetailIconTableActions';

interface ProjectDetailIconTableProps {
  data: IconModel[];
  projectId: number;
  onEdit?: (icon: IconModel) => void;
  onDelete?: (icon: IconModel) => void;
  loading?: boolean;
}

const statusColors = {
  pending:
    'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300',
  in_progress: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
  published:
    'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
  rejected: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300',
};

export default function ProjectDetailIconTable({
  data,
  projectId,
  onEdit,
  onDelete,
  loading = false,
}: ProjectDetailIconTableProps) {
  const [sorting, setSorting] = React.useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    []
  );
  const [globalFilter, setGlobalFilter] = React.useState('');

  const columns: ColumnDef<IconModel>[] = React.useMemo(
    () => [
      {
        id: 'actions',
        header: 'Actions',
        cell: ({ row }) => {
          const icon = row.original;
          return (
            <ProjectDetailIconTableActions
              icon={icon}
              onEdit={onEdit}
              onDelete={onDelete}
            />
          );
        },
      },
      {
        accessorKey: 'name',
        header: 'Name',
        cell: ({ row }) => (
          <div className='font-medium'>{row.getValue('name')}</div>
        ),
      },
      {
        accessorKey: 'pkg',
        header: 'Package',
        cell: ({ row }) => (
          <div className='font-mono text-sm text-muted-foreground'>
            {row.getValue('pkg')}
          </div>
        ),
      },
      {
        accessorKey: 'componentInfo',
        header: 'Component Info',
        cell: ({ row }) => (
          <div className='font-mono text-sm text-muted-foreground max-w-xs truncate'>
            {row.getValue('componentInfo')}
          </div>
        ),
      },
      {
        accessorKey: 'drawable',
        header: 'Drawable',
        cell: ({ row }) => {
          const icon = row.original;
          return (
            <div className='flex items-center space-x-2'>
              <ProjectDetailIconDrawableView
                projectId={projectId}
                drawable={icon.drawable}
                componentInfo={icon.componentInfo}
              />
              <div className='font-mono text-sm'>
                {row.getValue('drawable')}
              </div>
            </div>
          );
        },
      },
      {
        accessorKey: 'status',
        header: 'Status',
        cell: ({ row }) => {
          const status = row.getValue('status') as string;
          return (
            <Badge
              className={statusColors[status as keyof typeof statusColors]}
            >
              {status.replace('_', ' ')}
            </Badge>
          );
        },
      },
      {
        accessorKey: 'createdAt',
        header: 'Created',
        cell: ({ row }) => (
          <div className='text-sm text-muted-foreground'>
            {new Date(row.getValue('createdAt')).toLocaleDateString()}
          </div>
        ),
      },
    ],
    [onEdit, onDelete, projectId]
  );

  const table = useReactTable({
    data,
    columns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onGlobalFilterChange: setGlobalFilter,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    state: {
      sorting,
      columnFilters,
      globalFilter,
    },
  });

  return (
    <div className='space-y-4 w-full'>
      <div className='flex items-center justify-end'>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant='outline'>
              Columns <ChevronDown className='ml-2 h-4 w-4' />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align='end'>
            {table
              .getAllColumns()
              .filter(column => column.getCanHide())
              .map(column => {
                return (
                  <DropdownMenuCheckboxItem
                    key={column.id}
                    className='capitalize'
                    checked={column.getIsVisible()}
                    onCheckedChange={value => column.toggleVisibility(!!value)}
                  >
                    {column.id}
                  </DropdownMenuCheckboxItem>
                );
              })}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <div className='rounded-md border w-full overflow-x-auto'>
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map(headerGroup => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map(header => {
                  return (
                    <TableHead
                      key={header.id}
                      className='text-center px-4 py-3'
                    >
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className='h-24 text-center'
                >
                  <div className='flex items-center justify-center'>
                    <LoadingSpinner size={20} />
                  </div>
                </TableCell>
              </TableRow>
            ) : table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map(row => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && 'selected'}
                >
                  {row.getVisibleCells().map(cell => (
                    <TableCell key={cell.id} className='px-4 py-3'>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className='h-24 text-center'
                >
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <div className='flex items-center justify-end space-x-2 py-4'>
        <div className='flex-1 text-sm text-muted-foreground'>
          {table.getFilteredSelectedRowModel().rows.length} of{' '}
          {table.getFilteredRowModel().rows.length} row(s) selected.
        </div>
      </div>
    </div>
  );
}
