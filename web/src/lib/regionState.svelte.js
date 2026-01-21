import { browser } from '$app/environment';

class RegionState {
    value = $state(browser ? localStorage.getItem('region') || 'americas' : 'americas');

    label = $derived.by(() => {
        switch (this.value) {
            case 'americas':
                return 'Americas';
            case 'europe':
                return 'Europe';
            case 'asia':
                return 'Asia';
            default:
                return 'Americas';
        }
    });

    constructor() {
        $effect.root(() => {
            $effect(() => {
                if (browser) {
                    localStorage.setItem('region', this.value);
                }
            });
        });
    }
}

export const regionState = new RegionState();
