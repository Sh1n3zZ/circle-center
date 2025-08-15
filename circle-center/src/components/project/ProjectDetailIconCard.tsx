import { iconApi } from '@/api/manager/icon';
import type {
  IconModel,
  IconStats,
  ListIconsParams,
} from '@/api/manager/types';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Plus, RefreshCw } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import ProjectDetailIconEdit from './ProjectDetailIconEdit';
import ProjectDetailIconTable from './ProjectDetailIconTable';

interface ProjectDetailIconCardProps {
  projectId: number;
  onAddIcon?: () => void;
  onEditIcon?: (icon: IconModel) => void;
  onDeleteIcon?: (icon: IconModel) => void;
}

export default function ProjectDetailIconCard({
  projectId,
  onAddIcon,
  onEditIcon,
  onDeleteIcon,
}: ProjectDetailIconCardProps) {
  const [icons, setIcons] = useState<IconModel[]>([]);
  const [stats, setStats] = useState<IconStats | null>(null);
  const [loading, setLoading] = useState(false);
  const [filters, setFilters] = useState<ListIconsParams>({
    limit: 50,
    offset: 0,
  });
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);

  // Edit dialog state
  const [editingIcon, setEditingIcon] = useState<IconModel | null>(null);
  const [editDialogOpen, setEditDialogOpen] = useState(false);

  const fetchIcons = async () => {
    try {
      setLoading(true);
      const response = await iconApi.list(projectId, filters);
      setIcons(response.data.icons);
      setTotal(response.data.total);
      setTotalPages(response.data.totalPages);
      setCurrentPage(response.data.currentPage);
    } catch (error: any) {
      toast.error(error?.response?.data?.message || 'Failed to load icons');
    } finally {
      setLoading(false);
    }
  };

  const fetchStats = async () => {
    try {
      const response = await iconApi.getStats(projectId);
      setStats(response.data);
    } catch (error: any) {
      console.error('Failed to load stats:', error);
    }
  };

  useEffect(() => {
    fetchIcons();
    fetchStats();
  }, [projectId, filters]);

  const handleDelete = async (icon: IconModel) => {
    if (!onDeleteIcon) return;

    try {
      await iconApi.delete(projectId, icon.id);
      toast.success('Icon deleted successfully');
      fetchIcons();
      fetchStats();
    } catch (error: any) {
      toast.error(error?.response?.data?.message || 'Failed to delete icon');
    }
  };

  const handleEdit = (icon: IconModel) => {
    setEditingIcon(icon);
    setEditDialogOpen(true);
  };

  const handleEditSave = (updatedIcon: IconModel) => {
    // Update the icon in the local state
    setIcons(prev =>
      prev.map(icon => (icon.id === updatedIcon.id ? updatedIcon : icon))
    );

    // Call the parent callback if provided
    onEditIcon?.(updatedIcon);

    // Refresh data
    fetchIcons();
    fetchStats();
  };

  const handleStatusFilter = (status: string) => {
    setFilters(prev => ({
      ...prev,
      status: status === 'all' ? undefined : status,
      offset: 0, // Reset to first page when filtering
    }));
  };

  const handleSearch = (search: string) => {
    setFilters(prev => ({
      ...prev,
      search: search || undefined,
      offset: 0, // Reset to first page when searching
    }));
  };

  const handlePageChange = (newOffset: number) => {
    setFilters(prev => ({
      ...prev,
      offset: newOffset,
    }));
  };

  const handlePageJump = (page: number) => {
    const newOffset = (page - 1) * (filters.limit || 50);
    setFilters(prev => ({
      ...prev,
      offset: newOffset,
    }));
  };

  return (
    <>
      <Card>
        <CardHeader>
          <div className='flex items-center justify-between'>
            <CardTitle className='flex items-center gap-2'>
              Icons
              {stats && (
                <Badge variant='secondary'>{stats.total_icons} total</Badge>
              )}
            </CardTitle>
            <div className='flex items-center gap-2'>
              <Button
                variant='outline'
                size='sm'
                onClick={() => {
                  fetchIcons();
                  fetchStats();
                }}
                disabled={loading}
              >
                <RefreshCw className='h-4 w-4' />
              </Button>
              {onAddIcon && (
                <Button size='sm' onClick={onAddIcon}>
                  <Plus className='h-4 w-4 mr-2' />
                  Add Icon
                </Button>
              )}
            </div>
          </div>

          {/* Stats Row */}
          {stats && (
            <div className='flex items-center gap-4 text-sm'>
              <div className='flex items-center gap-1'>
                <Badge
                  variant='outline'
                  className='bg-yellow-50 text-yellow-700'
                >
                  {stats.pending_count}
                </Badge>
                <span className='text-muted-foreground'>Pending</span>
              </div>
              <div className='flex items-center gap-1'>
                <Badge variant='outline' className='bg-blue-50 text-blue-700'>
                  {stats.in_progress_count}
                </Badge>
                <span className='text-muted-foreground'>In Progress</span>
              </div>
              <div className='flex items-center gap-1'>
                <Badge variant='outline' className='bg-green-50 text-green-700'>
                  {stats.published_count}
                </Badge>
                <span className='text-muted-foreground'>Published</span>
              </div>
              <div className='flex items-center gap-1'>
                <Badge variant='outline' className='bg-red-50 text-red-700'>
                  {stats.rejected_count}
                </Badge>
                <span className='text-muted-foreground'>Rejected</span>
              </div>
            </div>
          )}
        </CardHeader>

        <CardContent className='w-full'>
          {/* Filters */}
          <div className='flex items-center gap-4 mb-4'>
            <div className='flex-1'>
              <Input
                placeholder='Search icons...'
                onChange={e => handleSearch(e.target.value)}
                className='max-w-sm'
              />
            </div>
            <Select onValueChange={handleStatusFilter} defaultValue='all'>
              <SelectTrigger className='w-40'>
                <SelectValue placeholder='All Status' />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value='all'>All Status</SelectItem>
                <SelectItem value='pending'>Pending</SelectItem>
                <SelectItem value='in_progress'>In Progress</SelectItem>
                <SelectItem value='published'>Published</SelectItem>
                <SelectItem value='rejected'>Rejected</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <ProjectDetailIconTable
            data={icons}
            projectId={projectId}
            onEdit={handleEdit}
            onDelete={handleDelete}
            loading={loading}
          />

          {total > 0 && (
            <div className='flex items-center justify-between mt-4'>
              <div className='text-sm text-muted-foreground'>
                Showing {(filters.offset || 0) + 1} to{' '}
                {Math.min((filters.offset || 0) + (filters.limit || 50), total)}{' '}
                of {total} icons
              </div>
              <div className='flex items-center gap-2'>
                <Button
                  variant='outline'
                  size='sm'
                  onClick={() =>
                    handlePageChange(
                      Math.max(0, (filters.offset || 0) - (filters.limit || 50))
                    )
                  }
                  disabled={(filters.offset || 0) === 0}
                >
                  Previous
                </Button>

                {/* Page numbers */}
                <div className='flex items-center gap-1'>
                  {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                    let pageNum;
                    if (totalPages <= 5) {
                      pageNum = i + 1;
                    } else if (currentPage <= 3) {
                      pageNum = i + 1;
                    } else if (currentPage >= totalPages - 2) {
                      pageNum = totalPages - 4 + i;
                    } else {
                      pageNum = currentPage - 2 + i;
                    }

                    return (
                      <Button
                        key={pageNum}
                        variant={
                          currentPage === pageNum ? 'default' : 'outline'
                        }
                        size='sm'
                        className='w-8 h-8 p-0'
                        onClick={() => handlePageJump(pageNum)}
                      >
                        {pageNum}
                      </Button>
                    );
                  })}

                  {totalPages > 5 && currentPage < totalPages - 2 && (
                    <>
                      {currentPage < totalPages - 3 && (
                        <span className='px-2'>...</span>
                      )}
                      <Button
                        variant='outline'
                        size='sm'
                        className='w-8 h-8 p-0'
                        onClick={() => handlePageJump(totalPages)}
                      >
                        {totalPages}
                      </Button>
                    </>
                  )}
                </div>

                <Button
                  variant='outline'
                  size='sm'
                  onClick={() =>
                    handlePageChange(
                      (filters.offset || 0) + (filters.limit || 50)
                    )
                  }
                  disabled={
                    (filters.offset || 0) + (filters.limit || 50) >= total
                  }
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Edit Icon Dialog */}
      <ProjectDetailIconEdit
        projectId={projectId}
        icon={editingIcon}
        open={editDialogOpen}
        onOpenChange={setEditDialogOpen}
        onSave={handleEditSave}
      />
    </>
  );
}
