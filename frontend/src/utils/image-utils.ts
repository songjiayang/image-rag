/**
 * Image utility functions for handling common image-related operations
 */

/**
 * Handle image loading errors by setting a placeholder image
 * Prevents infinite loops when placeholder also fails to load
 */
export const handleImageError = (event: Event, placeholderPath = '/images/placeholder-image.jpg') => {
  const img = event.target as HTMLImageElement
  if (!img) return

  // Prevent infinite loop by checking if already set to placeholder
  if (img.src !== placeholderPath) {
    img.src = placeholderPath
  } else {
    // If placeholder also fails, hide the image to stop requests
    img.style.display = 'none'
    console.warn('Placeholder image failed to load, hiding image element')
  }
}

/**
 * Format file size to human readable format
 */
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * Check if a file is a valid image
 */
export const isValidImage = (file: File): boolean => {
  return file.type.startsWith('image/')
}

/**
 * Get image dimensions from a file
 */
export const getImageDimensions = (file: File): Promise<{ width: number; height: number }> => {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => {
      resolve({ width: img.naturalWidth, height: img.naturalHeight })
    }
    img.onerror = reject
    img.src = URL.createObjectURL(file)
  })
}

/**
 * Convert data URL to File object
 */
export const dataURLToFile = (dataURL: string, filename: string): File => {
  const arr = dataURL.split(',')
  const mime = arr[0].match(/:(.*?);/)?.[1] || 'image/jpeg'
  const bstr = atob(arr[1])
  let n = bstr.length
  const u8arr = new Uint8Array(n)
  while (n--) {
    u8arr[n] = bstr.charCodeAt(n)
  }
  return new File([u8arr], filename, { type: mime })
}

/**
 * Create object URL from file and revoke it when done
 */
export const createTemporaryUrl = (file: File): string => {
  return URL.createObjectURL(file)
}

/**
 * Revoke object URL to free memory
 */
export const revokeTemporaryUrl = (url: string): void => {
  URL.revokeObjectURL(url)
}