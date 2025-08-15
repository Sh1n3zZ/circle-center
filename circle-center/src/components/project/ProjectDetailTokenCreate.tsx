import { tokenApi } from '@/api/manager/token';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Dices } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';

export default function ProjectDetailTokenCreate({
  projectId,
  onCreated,
}: {
  projectId: number;
  onCreated?: () => void;
}) {
  const [name, setName] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [plainToken, setPlainToken] = useState('');

  const randomizeName = () => {
    // Simple human-friendly random name
    const rnd = Math.random().toString(36).slice(2, 8);
    setName(`token-${rnd}`);
    toast.success('Random name generated');
  };

  const copyToken = async () => {
    try {
      await navigator.clipboard.writeText(plainToken);
      toast.success('Token copied');
    } catch {
      toast.error('Copy failed');
    }
  };

  const handleCreate = async () => {
    try {
      setSubmitting(true);
      const res = await tokenApi.create(projectId, {
        name: name.trim() || undefined,
      });
      setPlainToken(res.data.token);
      setDialogOpen(true);
      setName('');
      onCreated?.();
    } catch (e: any) {
      toast.error(e?.response?.data?.message || e.message || 'Create failed');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <div className='space-y-3'>
        <div className='relative'>
          <Input
            placeholder='Token name (optional)'
            value={name}
            onChange={e => setName(e.target.value)}
            className='pr-10'
          />
          <Button
            type='button'
            variant='ghost'
            size='sm'
            className='absolute right-0 top-0 h-full px-3 hover:bg-muted hover:text-primary transition-colors cursor-pointer'
            onClick={randomizeName}
          >
            <Dices className='h-4 w-4' />
          </Button>
        </div>
        <Button onClick={handleCreate} disabled={submitting} className='w-full'>
          {submitting ? 'Creating...' : 'Create Token'}
        </Button>
      </div>

      <AlertDialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Token created</AlertDialogTitle>
            <AlertDialogDescription>
              Please store this token securely. It will not be shown again.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <div className='rounded-md border bg-muted px-3 py-2 text-sm font-mono break-all'>
            {plainToken || ''}
          </div>
          <AlertDialogFooter>
            <Button variant='outline' onClick={copyToken}>
              Copy
            </Button>
            <AlertDialogAction>Done</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
