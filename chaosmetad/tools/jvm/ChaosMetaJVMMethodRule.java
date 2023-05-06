import org.json.JSONObject;

public class ChaosMetaJVMMethodRule {
    String method;
    String fault;
    String content;
    int lineNum;
    String importPkg;

    public static final String InsertAtFault = "insertAt";
    public static final String InsertAfterFault = "insertAfter";
    public static final String InsertBeforeFault = "insertBefore";

    public static final String SetBodyFault = "setBody";

    public static final String MethodKey = "Method";
    public static final String FaultKey = "Fault";
    public static final String ContentKey = "Content";
    public static final String LineNumKey = "LineNum";

    public static final String ImportPkgKey = "ImportPkg";


    public ChaosMetaJVMMethodRule(JSONObject jsonObject) throws Exception {
        if (jsonObject == null) {
            throw new Exception("must not null");
        }

        if (!jsonObject.has(MethodKey) || !jsonObject.has(FaultKey) || !jsonObject.has(ContentKey)) {
            throw new Exception("one of method、fault、content is not provide");
        }

        method = jsonObject.getString(MethodKey);
        fault = jsonObject.getString(FaultKey);
        content = jsonObject.getString(ContentKey);

        if (method.isEmpty() || fault.isEmpty() || content.isEmpty()) {
            throw new Exception("one of method、fault、content is empty");
        }

        lineNum = 0;
        if (fault.equals(InsertAtFault)) {
            if (!jsonObject.has(LineNumKey)) {
                throw new Exception("must provide args \"LineNum\" in \"insertAt\" fault");
            }
            lineNum = jsonObject.getInt(LineNumKey);
            if (lineNum < 0) {
                throw new Exception("\"LineNum\" must >= 0");
            }
        }

        importPkg = "";
        if (jsonObject.has(ImportPkgKey)) {
            importPkg = jsonObject.getString(ImportPkgKey);
        }
    }
}