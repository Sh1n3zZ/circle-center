import { Link } from 'react-router-dom'
import {
  NavigationMenu,
  NavigationMenuList,
  NavigationMenuItem,
  NavigationMenuLink,
  navigationMenuTriggerStyle,
} from '@/components/ui/navigation-menu'
import { Button } from '@/components/ui/button'
import { UserProfileNavBubble } from '@/components/user/UserProfileNavBubble'
import { authHelpers } from '@/api/client'

const Header = () => {
  const isAuthenticated = authHelpers.isAuthenticated()
  const currentUser = authHelpers.getCurrentUser<{
    display_name?: string
    username?: string
    avatar_url?: string
  }>()

  const handleEditProfile = () => {
    // TODO: Navigate to profile edit page
    console.log('Edit profile clicked')
  }

  const handleLogout = async () => {
    try {
      await authHelpers.clearAuthData()
      // Redirect to login page or refresh the page
      window.location.href = '/login'
    } catch (error) {
      console.error('Logout failed:', error)
    }
  }

  return (
    <header className="sticky top-0 z-40 w-screen backdrop-blur-md bg-white/70 dark:bg-gray-900/70 border-b border-border">
      <div className="mx-auto flex max-w-7xl items-center justify-between px-4 md:px-8 py-3">
        {/* Logo or Brand */}
        <Link to="/" className="text-lg font-semibold select-none">
          Circle&nbsp;Center
        </Link>

        {/* Navigation */}
        <NavigationMenu viewport={false}>
          <NavigationMenuList>
            <NavigationMenuItem>
              <NavigationMenuLink
                asChild
                className={navigationMenuTriggerStyle()}
              >
                <Link to="/">Home</Link>
              </NavigationMenuLink>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <NavigationMenuLink
                asChild
                className={navigationMenuTriggerStyle()}
              >
                <Link to="/about">About</Link>
              </NavigationMenuLink>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <NavigationMenuLink
                asChild
                className={navigationMenuTriggerStyle()}
              >
                <Link to="/reader">Reader</Link>
              </NavigationMenuLink>
            </NavigationMenuItem>
          </NavigationMenuList>
        </NavigationMenu>

        {/* User Actions */}
        <div className="flex items-center gap-2">
          {isAuthenticated ? (
            <UserProfileNavBubble
              displayName={currentUser?.display_name || currentUser?.username}
              avatarPath={currentUser?.avatar_url}
              size={36}
              onEditProfile={handleEditProfile}
              onLogout={handleLogout}
            />
          ) : (
            <>
              <Button variant="outline" asChild>
                <Link to="/login">Login</Link>
              </Button>
              <Button asChild>
                <Link to="/login?tab=register">Register</Link>
              </Button>
            </>
          )}
        </div>
      </div>
    </header>
  )
}

export default Header
