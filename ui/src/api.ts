import axios from 'axios';
import type { Repo } from './types';
/**
 * api封装
 */
export class Api {
    private host: string = "/api/";
    /**
     * get app info
     * @returns 
     */
    public async getAppInfo():Promise<string> {
        const response = await axios.get(this.host + "app");
        return response.data.title;
    }
    /**
     * get repoes
     * @param keyword keyword
     * @param page page num, start from 1
     * @returns 
     */
    public async getRepoes(keyword:string, page:number):Promise<[Repo[], boolean]> {
        const response = await axios.get(this.host + "repoes", {
            params: {
                keyword: keyword,
                page: page
            },
        });
        if (response.data.error) {
            throw response.data.error;
        }
        return [response.data.data, response.data.has_more];
    }

    /**
     * get repo detail
     * @param name 
     * @param refresh 
     * @returns 
     */
    public async getRepoDetail(name:string, refresh:boolean):Promise<Repo> {
        const response = await axios.get(this.host + "repo", {
            params: {
                name: name,
                refresh: refresh ? 1 : 0
            },
        });
        if (response.data.error) {
            throw response.data.error;
        }
        return response.data.data;
    }

}