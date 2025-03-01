import { useUser } from '@src/UserContext.tsx'
import { FC } from 'react'
import GravatarImage from '@src/components/GravatarImage.tsx'

export const CurrentUser: FC<{ display: 'icon' | 'name' }> = ({ display }) => {
  const user = useUser()
  return (
    <>
      {display == 'icon' && user.profileImageURL && (
        <GravatarImage
          profileUrl={user.profileImageURL}
          baseSize={32}
          style={{ height: '2em', top: '-0.2em', position: 'relative' }}
        />
      )}
      {display == 'name' && <span>{user.displayName}</span>}
    </>
  )
}

export default CurrentUser
