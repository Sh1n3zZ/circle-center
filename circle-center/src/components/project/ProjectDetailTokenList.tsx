import { tokenApi } from '@/api/manager/token';
import type { ProjectTokenModel } from '@/api/manager/types';
import { Button } from '@/components/ui/button';
import { LoadingSpinner } from '@/components/ui/loading-spinner';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

export default function ProjectDetailTokenList({
  projectId,
  refreshSignal = 0,
}: {
  projectId: number;
  refreshSignal?: number;
}) {
  const [tokens, setTokens] = useState<ProjectTokenModel[]>([]);
  const [loading, setLoading] = useState(false);
  const [deletingIds, setDeletingIds] = useState<Set<number>>(new Set());

  const fetchTokens = async () => {
    try {
      setLoading(true);
      const res = await tokenApi.list(projectId);
      setTokens(res.data || []);
    } catch (e: any) {
      toast.error(
        e?.response?.data?.message || e.message || 'Failed to load tokens'
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTokens();
  }, [projectId, refreshSignal]);

  const handleDelete = async (id: number) => {
    try {
      setDeletingIds(prev => new Set(prev).add(id));
      await tokenApi.delete(projectId, id);
      toast.success('Token deleted');
      setTokens(prev => prev.filter(t => t.id !== id));
    } catch (e: any) {
      toast.error(e?.response?.data?.message || e.message || 'Delete failed');
    } finally {
      setDeletingIds(prev => {
        const newSet = new Set(prev);
        newSet.delete(id);
        return newSet;
      });
    }
  };

  return (
    <div className='space-y-3'>
      {loading ? (
        <div className='flex items-center justify-center py-8'>
          <LoadingSpinner size={24} />
        </div>
      ) : tokens.length === 0 ? (
        <div className='text-sm text-muted-foreground'>No tokens yet</div>
      ) : (
        <div className='space-y-2'>
          {tokens.map(t => {
            const isDeleting = deletingIds.has(t.id);
            return (
              <div
                key={t.id}
                className='border rounded-md p-3 flex items-center justify-between'
              >
                <div className='space-y-0.5'>
                  <div className='font-medium'>{t.name}</div>
                  <div className='text-xs text-muted-foreground'>
                    Active: {t.active ? 'Yes' : 'No'} · Created:{' '}
                    {new Date(t.created_at).toLocaleString()} · Last used:{' '}
                    {t.last_used_at
                      ? new Date(t.last_used_at).toLocaleString()
                      : 'Never'}
                  </div>
                </div>
                <div className='flex gap-2'>
                  <Button
                    variant='destructive'
                    size='sm'
                    onClick={() => handleDelete(t.id)}
                    disabled={isDeleting}
                    className='min-w-[80px]'
                  >
                    {isDeleting ? <LoadingSpinner size={14} /> : 'Delete'}
                  </Button>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
