import { Link } from 'react-router-dom'
import {
  NavigationMenu,
  NavigationMenuList,
  NavigationMenuItem,
  NavigationMenuLink,
  navigationMenuTriggerStyle,
} from '@/components/ui/navigation-menu'

const Header = () => {
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
          </NavigationMenuList>
        </NavigationMenu>
      </div>
    </header>
  )
}

export default Header
