import org.json.JSONArray;
import org.json.JSONObject;
import org.json.JSONTokener;

import java.io.*;
import java.lang.instrument.ClassFileTransformer;
import java.lang.instrument.Instrumentation;
import java.lang.reflect.Constructor;
import java.time.Duration;
import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.Map;

public class ChaosMetaJVMAgent {

    private static final String InjectAction = "inject";
    private static final String RecoverAction = "recover";

    public static void agentmain(String args, Instrumentation inst) throws Exception {
        System.out.println("[info]agentmain load");

        long durationSecond;
        JSONObject ruleJson = getRuleJson(args);
        JSONArray classRuleList = ruleJson.getJSONArray("ClassList");
        durationSecond = ruleJson.getLong("Duration");
        if (durationSecond < 0) {
            throw new Exception("Duration must >= 0");
        }

        LocalDateTime startTime = LocalDateTime.now();

//        URLClassLoader classLoader = new URLClassLoader(new URL[]{new File(transformerJar).toURI().toURL()});
//        Class<?> transformerClass = classLoader.loadClass("ChaosMetaClassFileTransformer");
        Class<?> transformerClass = ChaosMetaClassFileTransformer.class;
        Constructor<?> constructor = transformerClass.getConstructor(String.class, ChaosMetaJVMMethodRule[].class);

        Map<String, Class> loadedClassesMap = new HashMap<>();
        Class[] allClass = inst.getAllLoadedClasses();
        for (Class c : allClass) {
            loadedClassesMap.put(c.getName(), c);
        }

        // do inject
//        boolean hasInjected;
        Map<Class, ClassFileTransformer> injectedTransformerMap = getTransformerMap(classRuleList, loadedClassesMap, constructor);
        try {
            doInject(inst, injectedTransformerMap);
//            hasInjected = true;
        } catch (Exception e) {
            String errMsg = String.format("[error]inject execute fail: %s", e.getMessage());
            System.out.println(errMsg);
            doRecover(inst, injectedTransformerMap);
//            hasInjected = false;
            throw new Exception(errMsg);
        }

        while (true) {
            Boolean ifTimeout = isTimeout(startTime, durationSecond);
            FileReader fileReader = null;
            try {
                fileReader = new FileReader(args);
            } catch (FileNotFoundException e) {
                System.out.println("[info]config file not exist, start to recover");
            }

            if (fileReader == null || ifTimeout) {
                doRecover(inst, injectedTransformerMap);
                if (fileReader != null) {
                    fileReader.close();
                }
                break;
            }

            fileReader.close();
            Thread.sleep(2000);
        }

//        classLoader.close();
        System.out.println("[info]agentmain finish");
    }

    private static Boolean isTimeout(LocalDateTime startTime, long durationSecond) {
        if (durationSecond == 0) {
            return false;
        }
        Duration duration = Duration.between(startTime, LocalDateTime.now());
        return duration.getSeconds() > durationSecond;
    }

    private static JSONObject getRuleJson(String configPath) throws IOException {
        FileReader fileReader = new FileReader(configPath);
        JSONObject ruleJson = new JSONObject(new JSONTokener(new BufferedReader(fileReader)));
        fileReader.close();
        return ruleJson;
    }

    private static Map<Class, ClassFileTransformer> getTransformerMap(JSONArray classRuleList, Map<String, Class> loadedClassesMap, Constructor<?> constructor) throws Exception {
        Map<Class, ClassFileTransformer> injectedTransformerMap = new HashMap<>();
        for (Object classRule : classRuleList) {
            JSONObject classRuleJSON = (JSONObject) classRule;
            String className = classRuleJSON.getString("Class");
            JSONArray methodList = classRuleJSON.getJSONArray("MethodList");
            if (className.isEmpty() || methodList.length() == 0) {
                throw new Exception("[error]must provide \"Class\" and \"MethodList\" of each rule");
            }

            ChaosMetaJVMMethodRule[] methodRules = new ChaosMetaJVMMethodRule[methodList.length()];
            for (int i = 0; i < methodList.length(); i++) {
                methodRules[i] = new ChaosMetaJVMMethodRule((JSONObject) methodList.get(i));
            }

            Class targetLoadedClass = loadedClassesMap.get(className);
            if (targetLoadedClass == null) {
                throw new Exception(String.format("[error]not found class %s in JVM", className));
            }

            ClassFileTransformer classFileTransformer = (ClassFileTransformer) constructor.newInstance(className, methodRules);
            injectedTransformerMap.put(targetLoadedClass, classFileTransformer);
        }

        return injectedTransformerMap;
    }

    private static void doInject(Instrumentation inst, Map<Class, ClassFileTransformer> injectedTransformerMap) throws Exception {
        System.out.println("[info]start to inject");
        for (Map.Entry<Class, ClassFileTransformer> entry : injectedTransformerMap.entrySet()) {
            inst.addTransformer(entry.getValue(), true);
            Class c = entry.getKey();
            try {
                executeRetrans(inst, c, (ChaosMetaClassFileTransformer) entry.getValue(), InjectAction);
            } catch (Exception e) {
                inst.removeTransformer(entry.getValue());
                throw new Exception(String.format("[error]inject class %s failed: %s", c.getName(), e.getMessage()));
            }

            System.out.printf("[info]inject class %s success\n", entry.getKey().getName());
        }
    }

    private static void doRecover(Instrumentation inst, Map<Class, ClassFileTransformer> injectedTransformerMap) {
        System.out.println("[info]start to recover");
        for (Map.Entry<Class, ClassFileTransformer> entry : injectedTransformerMap.entrySet()) {
            try {
                executeRetrans(inst, entry.getKey(), (ChaosMetaClassFileTransformer) entry.getValue(), RecoverAction);
                inst.retransformClasses(entry.getKey());
                System.out.printf("[info]recover class %s success\n", entry.getKey().getName());
            } catch (Exception e) {
                System.out.printf("[error]recover class %s failed: %s\n", entry.getKey().getName(), e.getMessage());
            }

            inst.removeTransformer(entry.getValue());
        }
    }

    // Because direct execution of retransformClasses cannot catch exceptions, this method is used
    private static void executeRetrans(Instrumentation inst, Class c, ChaosMetaClassFileTransformer transformer, String action) throws Exception {
        if (((action.equals(InjectAction) && transformer.getStatus().equals("initial"))) || ((action.equals(RecoverAction) && transformer.getStatus().equals("success")))) {
            inst.retransformClasses(c);
        } else {
            return;
        }

        if (transformer.getStatus().equals("fail")) {
            throw new Exception(transformer.getMessage());
        }
    }
}
