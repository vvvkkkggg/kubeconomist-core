import { Download, Search } from 'lucide-react';
import React from 'react';

interface ViewHeaderProps {
    onSearch: (query: string) => void;
    onHideEmptyToggle: (hidden: boolean) => void;
    onExport: () => void;
    showHideEmpty?: boolean;
    scanDateSelect?: React.ReactNode;
    onShowSavingsOnlyToggle?: (show: boolean) => void;
    showSavingsOnly?: boolean;
}

export const ViewHeader: React.FC<ViewHeaderProps> = ({
    onSearch,
    onHideEmptyToggle,
    onExport,
    showHideEmpty = true,
    scanDateSelect,
    onShowSavingsOnlyToggle,
    showSavingsOnly,
}) => {
    return (
        <div className="krr-header">
            <div className="controls-left">
                <div className="search-container">
                    <Search size={18} />
                    <input type="text" placeholder="Type to search" onChange={(e) => onSearch(e.target.value)} />
                </div>
                {showHideEmpty && (
                    <label className="checkbox-container">
                        <input type="checkbox" onChange={(e) => onHideEmptyToggle(e.target.checked)} />
                        Hide empty recommendations
                    </label>
                )}
                {showSavingsOnly && onShowSavingsOnlyToggle && (
                     <label className="checkbox-container">
                        <input type="checkbox" defaultChecked onChange={(e) => onShowSavingsOnlyToggle(e.target.checked)} />
                        Show savings only
                    </label>
                )}
            </div>
            <div className="controls-right">
                {scanDateSelect}
                <button className="export-button" onClick={onExport}>
                    <Download size={16} />
                    Export CSV
                </button>
            </div>
        </div>
    );
}; 
