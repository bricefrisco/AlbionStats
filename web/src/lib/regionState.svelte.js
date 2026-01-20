import { browser } from '$app/environment';

class RegionState {
    value = $state(browser ? localStorage.getItem('region') || 'americas' : 'americas');

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
