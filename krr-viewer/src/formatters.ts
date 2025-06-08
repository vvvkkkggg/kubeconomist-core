export function formatCPU(value: number | null | string): string {
    if (value === null || typeof value === 'string') {
        return value === '?' ? '?' : 'none';
    }
    if (value < 1) {
        return `${Math.round(value * 1000)}m`;
    }
    return `${value}`;
}

export function formatMemory(value: number | null | string): string {
    if (value === null || typeof value === 'string') {
        return value === '?' ? '?' : 'none';
    }
    const megaBytes = value / (1024 * 1024);
    if (megaBytes < 1024) {
        return `${Math.round(megaBytes)}Mi`;
    }
    const gigaBytes = megaBytes / 1024;
    return `${Math.round(gigaBytes * 100) / 100}Gi`;
}

export function formatChange(current: number | null | string, recommended: number | null | string, formatter: (val: number | null | string) => string): string {
    const currentFormatted = formatter(current);
    const recommendedFormatted = formatter(recommended);

    if (currentFormatted === '?' || recommendedFormatted === '?') {
        if (currentFormatted === 'none') {
            return `${currentFormatted} → ?`;
        }
        return `?`;
    }
    
    return `${currentFormatted} → ${recommendedFormatted}`;
} 
