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

// Calculate cost savings in rubles for container resource optimization
export function calculateCostSavings(scan: {
    object?: {
        allocations?: {
            requests?: {
                cpu?: number | null;
                memory?: number | null;
            };
        };
    };
    recommended?: {
        requests?: {
            cpu?: { value?: number | string | null };
            memory?: { value?: number | string | null };
        };
    };
} | null): number {
    if (!scan || !scan.object || !scan.recommended) {
        return 0;
    }

    // Hardcoded pricing from Yandex Cloud (standard-v1, 100% core fraction)
    // These should ideally come from a config or API, but matching the backend hardcoding
    const CPU_PRICE_PER_HOUR = 1.5; // rubles per core per hour
    const MEMORY_PRICE_PER_HOUR = 0.5; // rubles per GB per hour

    let savings = 0;

    // Calculate CPU savings
    const currentCPU = scan.object.allocations?.requests?.cpu;
    const recommendedCPU = scan.recommended.requests?.cpu?.value;
    
    if (currentCPU && recommendedCPU && 
        typeof currentCPU === 'number' && typeof recommendedCPU === 'number' && 
        currentCPU > recommendedCPU) {
        const cpuSavings = (currentCPU - recommendedCPU) * CPU_PRICE_PER_HOUR;
        savings += cpuSavings;
    }

    // Calculate Memory savings (convert from bytes to GB)
    const currentMemory = scan.object.allocations?.requests?.memory;
    const recommendedMemory = scan.recommended.requests?.memory?.value;
    
    if (currentMemory && recommendedMemory && 
        typeof currentMemory === 'number' && typeof recommendedMemory === 'number' && 
        currentMemory > recommendedMemory) {
        const currentMemoryGB = currentMemory / (1024 * 1024 * 1024);
        const recommendedMemoryGB = recommendedMemory / (1024 * 1024 * 1024);
        const memorySavings = (currentMemoryGB - recommendedMemoryGB) * MEMORY_PRICE_PER_HOUR;
        savings += memorySavings;
    }

    // Convert hourly savings to monthly (24 hours * 30 days = 720 hours per month)
    const monthlySavings = savings * 24 * 30;
    
    // Return monthly savings rounded to 2 decimal places
    return Math.round(monthlySavings * 100) / 100;
}

export function formatRubles(value: number): string {
    return `${value.toFixed(2)} ₽/мес`;
} 
