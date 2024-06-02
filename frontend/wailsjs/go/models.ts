export namespace publicCode {
	
	export class OtherSetting {
	    localLoading: boolean;
	
	    static createFrom(source: any = {}) {
	        return new OtherSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.localLoading = source["localLoading"];
	    }
	}
	export class ScanSetting {
	    timeout: number;
	    threadNum: number;
	
	    static createFrom(source: any = {}) {
	        return new ScanSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timeout = source["timeout"];
	        this.threadNum = source["threadNum"];
	    }
	}
	export class ProxySetting {
	    host: string;
	    port: string;
	    username: string;
	    password: string;
	    enable: boolean;
	    mode: string;
	
	    static createFrom(source: any = {}) {
	        return new ProxySetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.enable = source["enable"];
	        this.mode = source["mode"];
	    }
	}
	export class DBSetting {
	    mode: boolean;
	    host: string;
	    port: string;
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new DBSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mode = source["mode"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}
	export class AllSetting {
	    db: DBSetting;
	    proxy: ProxySetting;
	    scan: ScanSetting;
	    other: OtherSetting;
	
	    static createFrom(source: any = {}) {
	        return new AllSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.db = this.convertValues(source["db"], DBSetting);
	        this.proxy = this.convertValues(source["proxy"], ProxySetting);
	        this.scan = this.convertValues(source["scan"], ScanSetting);
	        this.other = this.convertValues(source["other"], OtherSetting);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Background {
	    image: string;
	
	    static createFrom(source: any = {}) {
	        return new Background(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.image = source["image"];
	    }
	}
	
	export class PocsInfoData {
	    uuid: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new PocsInfoData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.name = source["name"];
	    }
	}
	export class FingerprintScanData {
	    path: string;
	    request_method: string;
	    request_headers: {[key: string]: string};
	    request_data: string;
	    status_code: number;
	    headers: {[key: string]: string};
	    keyword: string[];
	    favicon_hash: string[];
	    priority: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new FingerprintScanData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.request_method = source["request_method"];
	        this.request_headers = source["request_headers"];
	        this.request_data = source["request_data"];
	        this.status_code = source["status_code"];
	        this.headers = source["headers"];
	        this.keyword = source["keyword"];
	        this.favicon_hash = source["favicon_hash"];
	        this.priority = source["priority"];
	        this.name = source["name"];
	    }
	}
	export class Fingerprint {
	    uuid: string;
	    name: string;
	    fingerprintScan: FingerprintScanData[];
	    pocsInfo: PocsInfoData[];
	
	    static createFrom(source: any = {}) {
	        return new Fingerprint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.name = source["name"];
	        this.fingerprintScan = this.convertValues(source["fingerprintScan"], FingerprintScanData);
	        this.pocsInfo = this.convertValues(source["pocsInfo"], PocsInfoData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class FingerprintScanResult {
	    url: string;
	    title: string;
	    fingerprint: string[];
	    vulnerability: string[];
	
	    static createFrom(source: any = {}) {
	        return new FingerprintScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.title = source["title"];
	        this.fingerprint = source["fingerprint"];
	        this.vulnerability = source["vulnerability"];
	    }
	}
	export class Need {
	    label: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new Need(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.value = source["value"];
	    }
	}
	
	export class RequestData {
	    pocString: string;
	    status: string;
	    check: string;
	    print: string;
	
	    static createFrom(source: any = {}) {
	        return new RequestData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pocString = source["pocString"];
	        this.status = source["status"];
	        this.check = source["check"];
	        this.print = source["print"];
	    }
	}
	export class Poc {
	    uuid: string;
	    name: string;
	    hunter: string;
	    fofa: string;
	    cms: string;
	    description: string;
	    optionValue: number;
	    needData: Need[];
	    request: RequestData[];
	
	    static createFrom(source: any = {}) {
	        return new Poc(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.name = source["name"];
	        this.hunter = source["hunter"];
	        this.fofa = source["fofa"];
	        this.cms = source["cms"];
	        this.description = source["description"];
	        this.optionValue = source["optionValue"];
	        this.needData = this.convertValues(source["needData"], Need);
	        this.request = this.convertValues(source["request"], RequestData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PocScanResult {
	    url: string;
	    pocName: string;
	    print: string;
	
	    static createFrom(source: any = {}) {
	        return new PocScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.pocName = source["pocName"];
	        this.print = source["print"];
	    }
	}
	
	
	

}

