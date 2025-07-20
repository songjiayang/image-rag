/**
 * Date utility functions for consistent date formatting across the application
 */

/**
 * Format date to a standard readable format
 */
export const formatDate = (dateString: string | Date): string => {
  const date = typeof dateString === 'string' ? new Date(dateString) : dateString
  return date.toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

/**
 * Format date with time
 */
export const formatDateTime = (dateString: string | Date): string => {
  const date = typeof dateString === 'string' ? new Date(dateString) : dateString
  return date.toLocaleString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

/**
 * Format relative time (e.g., "2 hours ago")
 */
export const formatRelativeTime = (dateString: string | Date): string => {
  const date = typeof dateString === 'string' ? new Date(dateString) : dateString
  const now = new Date()
  const diffInMs = now.getTime() - date.getTime()
  const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24))

  if (diffInDays === 0) {
    const diffInHours = Math.floor(diffInMs / (1000 * 60 * 60))
    if (diffInHours === 0) {
      const diffInMinutes = Math.floor(diffInMs / (1000 * 60))
      return diffInMinutes === 0 ? 'Just now' : `${diffInMinutes} minute${diffInMinutes > 1 ? 's' : ''} ago`
    }
    return `${diffInHours} hour${diffInHours > 1 ? 's' : ''} ago`
  } else if (diffInDays === 1) {
    return 'Yesterday'
  } else if (diffInDays < 7) {
    return `${diffInDays} day${diffInDays > 1 ? 's' : ''} ago`
  } else {
    return formatDate(date)
  }
}

/**
 * Get ISO string for API calls
 */
export const toISOString = (date: Date): string => {
  return date.toISOString()
}

/**
 * Parse ISO string to Date
 */
export const parseISOString = (isoString: string): Date => {
  return new Date(isoString)
}