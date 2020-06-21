## Current issue
- Archive extraction takes too long
- Cache is not implemented as a cache it has a mixture of Cache and Image Processor
- Already processed images are not reused when gnol restarts
- No Logs
- No proper Error Handling

## Design Goal
- Clear Separation of Cache and Cache filling
- Make Cache Recoverable, so that after a restart cache is recreated
- Create method to delete cach or parts of it with Last Used or Least frequently used stats
- Export Cache statistics
- Create a separate module to create images
- Create a separate module to get images from Archive
- Keep memory profile low
- Fast startup Times


## Implementation

  
