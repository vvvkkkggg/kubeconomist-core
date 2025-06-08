import type { PlatformOptimizerRecommendation } from '../types';

export const mockPlatformOptimizerRecommendations: PlatformOptimizerRecommendation[] = [
  {
    id: 'cat1b2c3d4e5f6g7h8i9',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'fo1e3k9lmn4p5q6r7s8',
    nodeGroupId: 'cat1b2c3d4e5f6g7h8i9',
    currentPlatform: 'standard-v2',
    desiredPlatform: 'standard-v3',
    currentPrice: 250.75,
    desiredPrice: 210.50,
    savings: 40.25,
  },
  {
    id: 'cat9i8h7g6f5e4d3c2b1',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'fo2a3b4c5d6e7f8g9h0',
    nodeGroupId: 'cat9i8h7g6f5e4d3c2b1',
    currentPlatform: 'standard-v1',
    desiredPlatform: 'standard-v3',
    currentPrice: 350.00,
    desiredPrice: 210.50,
    savings: 139.50,
  },
  {
    id: 'catf0e1d2c3b4a59876',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'b1g8enqi42dal6t489ou',
    nodeGroupId: 'catf0e1d2c3b4a59876',
    currentPlatform: 'standard-v2',
    desiredPlatform: 'standard-v3',
    currentPrice: 501.50,
    desiredPrice: 421.00,
    savings: 80.50,
  },
  // This one is already on the cheapest platform
  {
    id: 'cat67895a4b3c2d1e0f',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'b1g8enqi42dal6t489ou',
    nodeGroupId: 'cat67895a4b3c2d1e0f',
    currentPlatform: 'standard-v3',
    desiredPlatform: 'standard-v3',
    currentPrice: 199.99,
    desiredPrice: 199.99,
    savings: 0,
  },
]; 
