import javassist.ClassPool;
import javassist.CtClass;
import javassist.CtMethod;

import java.lang.instrument.ClassFileTransformer;
import java.lang.instrument.IllegalClassFormatException;
import java.security.ProtectionDomain;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;

public class ChaosMetaClassFileTransformer implements ClassFileTransformer {
    private final String targetClassName;

    private final ChaosMetaJVMMethodRule[] methodRules;

    private String status;

    private String message;

    private byte[] oldBytes;

    public ChaosMetaClassFileTransformer(String targetClassName, ChaosMetaJVMMethodRule[] methodRules) {
        this.targetClassName = targetClassName;
        this.methodRules = methodRules;
        this.oldBytes = null;
        this.status = "initial";
    }

    private void doImportPkg(String pkgStr, ClassPool pool) {
        Iterator<String> importedPkg = pool.getImportedPackages();
        Map<String, Boolean> importedPkgMap = new HashMap<>();
        while (importedPkg.hasNext()) {
            String pkg = importedPkg.next();
            importedPkgMap.put(pkg, true);
        }

        String[] pkgList = pkgStr.split(",");
        for (String pkg : pkgList) {
            if (importedPkgMap.get(pkg) == null || !importedPkgMap.get(pkg)) {
                System.out.printf("import pkg: %s\n", pkg);
                pool.importPackage(pkg);
            }
        }
    }

    private void doInject(CtMethod ctMethod, ChaosMetaJVMMethodRule methodRule) throws Exception {
        System.out.printf("inject fault: %s\n", methodRule.fault);
        switch (methodRule.fault) {
            case ChaosMetaJVMMethodRule.InsertBeforeFault:
                ctMethod.insertBefore(methodRule.content);
                break;
            case ChaosMetaJVMMethodRule.InsertAfterFault:
                ctMethod.insertAfter(methodRule.content);
                break;
            case ChaosMetaJVMMethodRule.SetBodyFault:
                ctMethod.setBody(methodRule.content);
                break;
            case ChaosMetaJVMMethodRule.InsertAtFault:
                ctMethod.insertAt(methodRule.lineNum, methodRule.content);
                break;
            default:
                throw new Exception("not support fault: " + methodRule.fault);
        }
    }


    public String getMessage() {
        return message;
    }

    public String getStatus() {
        return status;
    }

    // TODO: It needs to be processed so that it will not be affected even if it is loaded by other class loaders
    @Override
    public byte[] transform(ClassLoader loader, String className, Class<?> classBeingRedefined,
                            ProtectionDomain protectionDomain, byte[] classfileBuffer) throws IllegalClassFormatException {
        System.out.printf("transform class: %s, targetClassName: %s, status: %s, loader: %s, \n", className, targetClassName, status, loader.toString());

        try {
            if (!targetClassName.equals(className) || (status.equals("recovered"))) {
                return null;
            }

            ClassPool pool = ClassPool.getDefault();
            CtClass ctClass = pool.get(targetClassName);
            if (status.equals("success") || status.equals("fail")) {
                // TODO: Need to consider how to clean up the added packages when restoring?
                System.out.printf("recover class: %s\n", className);
                if (ctClass.isFrozen()) {
                    ctClass.defrost();
                }
                byte[] toReturn = oldBytes;
                oldBytes = null;
                status = "recovered";
                return toReturn;
            }
//            oldBytes = ctClass.toBytecode(); // It cannot be used, because it is just a reference copy, and it will be saved if it is modified later
//            if (ctClass.isFrozen()) {
//                ctClass.defrost();
//            }
            oldBytes = classfileBuffer;
            System.out.printf("inject class: %s\n", className);
            for (ChaosMetaJVMMethodRule methodRule : methodRules) {
                System.out.printf("resolve rule: %s\n", methodRule.method);
                if (!methodRule.importPkg.isEmpty()) {
                    doImportPkg(methodRule.importPkg, pool);
                }

                CtMethod ctMethod = ctClass.getDeclaredMethod(methodRule.method);
                doInject(ctMethod, methodRule);
            }

            byte[] toReturn = ctClass.toBytecode();
            // Prevents secondary modification, so it can remain frozen
//            if (ctClass.isFrozen()) {
//                ctClass.defrost();
//            }
            ctClass.detach();
            System.out.println("transform finish");
            this.status = "success";

            return toReturn;
        } catch (Exception e) {
            this.status = "fail";
            message = e.getMessage();
            throw new IllegalClassFormatException(message);
        }
    }
}
